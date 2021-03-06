package datastore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/yssk22/go/cache"
	"github.com/yssk22/go/iterator/slice"
	"github.com/yssk22/go/retry"
	"github.com/yssk22/go/x/xerrors"
	"github.com/yssk22/go/x/xlog"
	"github.com/yssk22/go/x/xnet"
	"github.com/yssk22/go/x/xruntime"
	"github.com/yssk22/go/x/xtime"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

const fixtureLoggerKey = "github.com/yssk22/gcp/datastore.fixture"

var _floatRe = regexp.MustCompile("\\.0+$")

type emulator struct {
	process *os.Process
	dir     string
	port    int
}

func startEmulator(pkgName string) (*emulator, error) {
	var err error
	var port int
	var dir string
	defer func() {
		if err != nil {
			log.Printf("failed to start an emulator process: %v", err)
		} else {
			log.Printf("start an emulator process at %d (%s)", port, dir)
		}
	}()
	port, err = xnet.GetEphemeralPort()
	if err != nil {
		err = xerrors.Wrap(err, "cannot start an emulator - ephemeral port assignment failure")
		return nil, err
	}
	dir, err = ioutil.TempDir("", "gcp-datastore-emulator")
	if err != nil {
		err = xerrors.Wrap(err, "cannot start an emulator - cannot create a tempdir")
		return nil, err
	}
	dir = filepath.Join(dir, pkgName)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		err = xerrors.Wrap(err, "cannot start an emulator - cannot create a package dir")
		return nil, err
	}
	args := []string{
		"beta",
		"emulators",
		"datastore",
		"start",
		"--consistency=1.0",
		fmt.Sprintf("--data-dir=%s", dir),
		"--no-store-on-disk",
		fmt.Sprintf("--host-port=localhost:%d", port),
		"--project=testenvironment",
	}
	var stdout, stderr *bytes.Buffer
	cliStr := fmt.Sprintf("gcloud %s", strings.Join(args, " "))
	cmd := exec.Command("gcloud", args...)
	if outputEnvironmentLogs() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		stdout = &(bytes.Buffer{})
		stderr = &(bytes.Buffer{})
		cmd.Stdout = stdout
		cmd.Stderr = stderr
	}
	err = cmd.Start()
	if err != nil {
		err = xerrors.Wrap(err, "cannot start an emulator - failed to start `%s`", cliStr)
		return nil, err
	}
	const timeout = 60 * time.Second
	interval := retry.ConstBackoff(200 * time.Millisecond)
	until := retry.Until(time.Now().Add(timeout))
	err = retry.Do(context.Background(), func(_ context.Context) error {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", port))
		defer func() {
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
		}()
		return err
	}, interval, until)
	if err != nil {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		if !outputEnvironmentLogs() {
			log.Println("failed to run a datastore emulator")
			log.Println("[stdout]")
			log.Println(stdout.String())
			log.Println("[stderr]")
			log.Println(stderr.String())
		}
		err = xerrors.Wrap(err, "cannot start an emulator: timedout in %s", timeout)
		return nil, err
	}
	return &emulator{
		process: cmd.Process,
		dir:     dir,
		port:    port,
	}, nil
}

func (e *emulator) Shutdown() error {
	var err error
	defer func() {
		if err != nil {
			log.Printf("failed to shutdown an emulator process at %d: %v", e.port, err)
		} else {
			log.Printf("shutdown an emulator process at %d (%s)", e.port, e.dir)
		}
	}()
	defer os.RemoveAll(e.dir)
	proc := e.process
	if proc == nil {
		err = fmt.Errorf("no emulator process")
		return err
	}
	errc := make(chan error, 1)
	go func() {
		_, err = proc.Wait()
		errc <- err
	}()
	var resp *http.Response
	resp, err = http.Post(fmt.Sprintf("http://localhost:%d/shutdown", e.port), "text/html", nil)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		if kerr := proc.Kill(); kerr != nil {
			err = xerrors.Wrap(kerr, "cannot kill the emulator main process")
			return err
		}
		err = xerrors.Wrap(err, "failed to request a shutdown")
		return err
	}
	select {
	case <-time.After(15 * time.Second):
		err = fmt.Errorf("shutdown timed out")
	case err = <-errc:
	}
	return err
}

// TestEnv is a struct to provide a helper
type TestEnv struct {
	context  context.Context
	memcache cache.Cache
	emulator *emulator
}

// NewTestEnv returns a new TestEnv instance
func NewTestEnv() (*TestEnv, error) {
	var pkgName string
	stack := xruntime.CaptureStack(10)
	for _, frame := range stack {
		if frame.FunctionName == "TestMain" {
			pkgName = frame.PackageName
			break
		}
	}
	ctx := context.Background()
	emulator, err := startEmulator(pkgName)
	if err != nil {
		return nil, err
	}
	// client, err := datastore.NewClient(ctx, "testenvironment",
	// 	option.WithEndpoint(fmt.Sprintf("localhost:%d", emulator.port)),
	// 	option.WithoutAuthentication(),
	// 	option.WithGRPCDialOption(grpc.WithInsecure()),
	// )
	// if err != nil {
	// 	return nil, err
	// }
	return &TestEnv{
		context:  ctx,
		memcache: &cache.MemoryCache{},
		emulator: emulator,
	}, nil
}

// MustNewTestEnv is like MustNewTestEnv, but panic if an error occurs
func MustNewTestEnv() *TestEnv {
	te, err := NewTestEnv()
	xerrors.MustNil(err)
	return te
}

// NewClient returns a *datastore.Client that sends requests to the test environment emulator
func (te *TestEnv) NewClient() *Client {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "testenvironment",
		option.WithEndpoint(fmt.Sprintf("localhost:%d", te.emulator.port)),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithInsecure()),
	)
	xerrors.MustNil(err)
	return NewClientFromClient(context.Background(), client, Cache(te.memcache))
}

// GetCache returns a cache client
func (te *TestEnv) GetCache() cache.Cache {
	return te.memcache
}

// Shutdown shuts down the environment
func (te *TestEnv) Shutdown() error {
	return te.emulator.Shutdown()
}

// Reset resets the environment
func (te *TestEnv) Reset() error {
	ctx := context.Background()
	te.memcache.Clear(ctx)
	// datastore cleanup
	client := te.NewClient()
	defer client.Close()
	q := datastore.NewQuery("__namespace__").KeysOnly()
	namespaceKeys, err := client.inner.GetAll(ctx, q, nil)
	if err != nil {
		return xerrors.Wrap(err, "cannot query namespaces")
	}
	return slice.Parallel(namespaceKeys, func(i int, nsKey *datastore.Key) error {
		q := datastore.NewQuery("__kind__").KeysOnly().Namespace(nsKey.Name)
		kindKeys, err := client.inner.GetAll(ctx, q, nil)
		if err != nil {
			return xerrors.Wrap(err, "cannot find kind keys")
		}
		for _, kindKey := range kindKeys {
			q := datastore.NewQuery(kindKey.Name).KeysOnly().Namespace(nsKey.Name)
			entityKeys, err := client.inner.GetAll(ctx, q, nil)
			if err != nil {
				return xerrors.Wrap(err, "cannot find entity keys")
			}
			if err = client.DeleteMulti(ctx, entityKeys); err != nil {
				return xerrors.Wrap(err, "cannot delete entities")
			}
		}
		return nil
	})
}

// LoadFixture loads the fixture data from `path`
func (te *TestEnv) LoadFixture(path string) error {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return xerrors.Wrap(err, "could not load fixture file from %s", path)
	}
	var arr []map[string]interface{}
	if err = json.Unmarshal(buff, &arr); err != nil {
		return xerrors.Wrap(err, "could not load the json file from %q", path)
	}
	for _, v := range arr {
		if err := te.json2Datastore(nil, v); err != nil {
			return err
		}
	}
	return nil
}

// MustLoadFixture is like LoadFixture but panic if an error occurs
func (te *TestEnv) MustLoadFixture(path string) {
	xerrors.MustNil(te.LoadFixture(path))
}

type jsonSaver map[string]interface{}

func (js jsonSaver) Load(ps []datastore.Property) error {
	return nil
}

func (js jsonSaver) Save() ([]datastore.Property, error) {
	props := []datastore.Property{}
	for k, v := range map[string]interface{}(js) {
		if !strings.HasPrefix(k, "_") {
			for _, val := range json2Properties(k, v) {
				props = append(props, val)
			}
		}
	}
	return props, nil
}

func loadFile(path string, bindings interface{}) ([]byte, error) {
	t, err := template.New(filepath.Base(path)).ParseFiles(path)
	if err != nil {
		return nil, err
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, bindings)
	return buff.Bytes(), err
}

func json2Properties(k string, v interface{}) []datastore.Property {
	var propertyList []datastore.Property
	var value = reflect.ValueOf(v)

	switch value.Kind() {
	case reflect.String:
		p := datastore.Property{Name: k}
		s := v.(string)
		if strings.HasPrefix(s, "[]") {
			p.Value = []byte(strings.TrimPrefix(s, "[]"))
			p.NoIndex = true
			propertyList = append(propertyList, p)
		} else {
			if dt, err := xtime.Parse(fmt.Sprintf("%s", v)); err == nil {
				p.Value = dt
				propertyList = append(propertyList, p)
			} else if d, err := xtime.Parse(fmt.Sprintf("%sT00:00:00Z", v)); err == nil {
				p.Value = d
				propertyList = append(propertyList, p)
			} else {
				p.Value = s
				propertyList = append(propertyList, p)
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// reach here from FixtureFromMap since it can contain non floating number.
		var vv int64
		switch v.(type) {
		case int:
			vv = int64(v.(int))
		case int8:
			vv = int64(v.(int8))
		case int16:
			vv = int64(v.(int16))
		case int32:
			vv = int64(v.(int32))
		case int64:
			vv = v.(int64)
		}
		propertyList = append(propertyList, datastore.Property{
			Name:  k,
			Value: vv,
		})
	case reflect.Float32, reflect.Float64:
		str := []byte(fmt.Sprintf("%f", v))
		if _floatRe.Match(str) {
			// should be int.
			propertyList = append(propertyList, datastore.Property{
				Name:  k,
				Value: int64(v.(float64)),
			})
		} else {
			propertyList = append(propertyList, datastore.Property{
				Name:  k,
				Value: v,
			})
		}
	case reflect.Bool:
		propertyList = append(propertyList, datastore.Property{
			Name:  k,
			Value: v,
		})
	case reflect.Map:
		for k1, v1 := range v.(map[string]interface{}) {
			if !strings.HasPrefix(k1, "_") {
				for _, val := range json2Properties(k1, v1) {
					propertyList = append(propertyList, datastore.Property{
						Name:  fmt.Sprintf("%s.%s", k, val.Name),
						Value: val.Value,
					})
				}
			}
		}
	case reflect.Slice:
		propertyList = append(propertyList, datastore.Property{
			Name:  k,
			Value: value.Interface(),
		})
	default:
		break
	}
	return propertyList
}

func (te *TestEnv) json2Datastore(pkey *datastore.Key, data map[string]interface{}) error {
	ctx, logger := xlog.WithContextAndKey(te.context, "", fixtureLoggerKey)
	client := te.NewClient()
	defer client.Close()

	var kind string
	var keyval interface{}
	var key *datastore.Key
	var ok bool
	if _, ok = data["_kind"]; !ok {
		return fmt.Errorf("missing key `_kind`")
	}
	kind = data["_kind"].(string)
	if keyval, ok = data["_key"]; !ok {
		return fmt.Errorf("missing key `_key`")
	}

	switch keyval.(type) {
	case string:
		key = datastore.NameKey(kind, keyval.(string), pkey)
		if _, ok = data["_ns"]; ok {
			key.Namespace = data["_ns"].(string)
		}
	default:
		return fmt.Errorf("invalid `_key` type for %v", keyval)
	}

	if _, err := client.inner.Put(ctx, key, jsonSaver(data)); err != nil {
		return err
	}
	if outputEnvironmentLogs() {
		if key.Namespace == "" {
			logger.Infof("Fixture: %s loaded", key)
		} else {
			logger.Infof("Fixture: %s%s loaded", key.Namespace, key)
		}
	}
	if children, ok := data["_children"]; ok {
		for _, v := range children.([]interface{}) {
			if err := te.json2Datastore(key, v.(map[string]interface{})); err != nil {
				return err
			}
		}
	}
	return nil
}

func outputEnvironmentLogs() bool {
	return os.Getenv("OUTPUT_TEST_ENVIRONMENT_LOG") == "1"
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}
