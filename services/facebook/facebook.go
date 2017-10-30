package facebook

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/speedland/go/x/xerrors"
	"github.com/speedland/go/x/xstrings"

	"github.com/speedland/go/x/xtime"

	"io/ioutil"

	"bytes"

	"context"

	"github.com/speedland/go/x/xlog"
)

const (
	domainGraphAPI      = "graph.facebook.com"
	domainGraphVideoAPI = "graph-video.facebook.com"
)

// LoggerKey is a key for this package
const LoggerKey = "speedland.net.services.facebook"

// Client is a API client for facebook graph API.
type Client struct {
	client      *http.Client
	accessToken string
	version     string
}

// NewClient returns a new facebook client
func NewClient(client *http.Client, accessToken string) *Client {
	return &Client{
		client:      client,
		accessToken: accessToken,
		version:     "v2.10",
	}
}

// SetVersion specifies the version used for the client.
func (c *Client) SetVersion(v string) {
	c.version = v
}

// Me is an struct returned by /me endpoint.
type Me struct {
	ID string `json:"id"`
}

// GetMe gets the *Profile of an authorized user.
func (c *Client) GetMe(ctx context.Context) (*Me, error) {
	var m Me
	if err := c.Get(ctx, "/me", url.Values{
		"fields": []string{"id"},
	}, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// URLObject is a object returned by scraping URL
type URLObject struct {
	ID          string           `json:"id"`
	URL         string           `json:"url"`
	Type        string           `json:"type"`
	Title       string           `json:"title"`
	UpdatedTime *xtime.Timestamp `json:"updated_time"`
	Image       []struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"image,omitempty"`
	Pages []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pages,omitempty"`
	Application struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"application,omitempty"`
}

func (c *Client) ScrapeURL(ctx context.Context, urlstr string) (*URLObject, error) {
	var obj URLObject
	err := c.Post(ctx, "/", url.Values{
		"id":          []string{urlstr},
		"scrape":      []string{"true"},
		"date_format": []string{"U"},
	}, nil, &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

const (
	errCouldNotGetResponse string = "could not get response from %s: %s"
	errInvalidStatusCode          = "could not get response from %s: got %d, expected: %d"
	errReadingResponse            = "could not read response from %s: %s"
	errInvalidJSONBody            = "could not parse json from %s: %s"
)

// Get to call GET request on the given path
func (c *Client) Get(ctx context.Context, path string, query url.Values, v interface{}) error {
	return c.get(ctx, domainGraphAPI, path, query, v)
}

// GetVideo to call GET request on the given path for video endpoints
func (c *Client) GetVideo(ctx context.Context, path string, query url.Values, v interface{}) error {
	return c.get(ctx, domainGraphVideoAPI, path, query, v)
}

func (c *Client) get(ctx context.Context, domain string, path string, query url.Values, v interface{}) error {
	url := c.prepareURL(domain, path, query, c.accessToken)
	_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("[GraphAPI:GET %s] ", url), LoggerKey)
	resp, err := c.client.Get(url)
	return processResponse(logger, v, url, nil, resp, err)
}

// Post to call POST request on the given path with json body specified by `content` argument.
func (c *Client) Post(ctx context.Context, path string, query url.Values, content interface{}, v interface{}) error {
	return c.post(ctx, c.prepareURL(domainGraphAPI, path, query, c.accessToken), content, v)
}

// MultipartParams is a struct to upload multipart/form data to the endpoint
type MultipartParams struct {
	ContentType string
	Body        io.Reader
}

// PostVideo to call POST request on the given path with json body specified by `content` argument.
func (c *Client) PostVideo(ctx context.Context, path string, query url.Values, content interface{}, v interface{}) error {
	return c.post(ctx, c.prepareURL(domainGraphVideoAPI, path, query, c.accessToken), content, v)
}

func (c *Client) post(ctx context.Context, urlstr string, content interface{}, v interface{}) error {
	_, logger := xlog.WithContextAndKey(ctx, fmt.Sprintf("[GraphAPI:POST %s] ", urlstr), LoggerKey)
	var buff bytes.Buffer
	var resp *http.Response
	var err error
	switch t := content.(type) {
	// TODO: case url.Values: for www-form-urlencoding
	case *MultipartParams:
		var req *http.Request
		params := content.(*MultipartParams)
		req, err = http.NewRequest("POST", urlstr, io.TeeReader(params.Body, &buff))
		req.Header.Set("Content-Type", params.ContentType)
		if err != nil {
			return xerrors.Wrap(err, "could not create a *http.NewRequest")
		}
		resp, err = c.client.Do(req)
	default:
		var jsonbuff bytes.Buffer
		if err = json.NewEncoder(&jsonbuff).Encode(content); err != nil {
			return xerrors.Wrap(err, "could not encode %s to JSON", t)
		}
		resp, err = c.client.Post(urlstr, "application/json; charset=utf-8", io.TeeReader(&jsonbuff, &buff))
	}
	return processResponse(logger, v, urlstr, &buff, resp, err)
}

func (c *Client) prepareURL(domain string, path string, query url.Values, accessToken string) string {
	if query == nil {
		query = url.Values{}
	}
	query.Set("access_token", accessToken)
	return fmt.Sprintf("https://%s/%s%s?%s", domain, c.version, path, query.Encode())
}

func processResponse(logger *xlog.Logger, v interface{}, url string, content interface{}, resp *http.Response, err error) error {
	if err != nil {
		return logAndError(
			logger,
			"",
			fmt.Errorf(errCouldNotGetResponse, url, err),
		)
	}
	defer resp.Body.Close()
	// expected json bytes are not so big so buffer here.
	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return logAndError(
			logger,
			"",
			fmt.Errorf(errReadingResponse, url, err),
		)
	}
	if resp.StatusCode != 200 {
		var message = fmt.Sprintf("response body: %s", string(buff))
		if resp.StatusCode == 400 || resp.StatusCode >= 500 && content != nil {
			if buff, ok := content.(io.Reader); ok {
				reqBody, _ := ioutil.ReadAll(buff)
				message = fmt.Sprintf("%s, request body (form): %s, request header: %s", message, string(reqBody), resp.Request.Header)
			} else {
				reqBody, _ := json.Marshal(content)
				message = fmt.Sprintf("%s, request body (json): %s", message, string(reqBody))
			}
		}
		return logAndError(
			logger,
			message,
			fmt.Errorf(errInvalidStatusCode, url, resp.StatusCode, 200),
		)
	}
	if err := json.Unmarshal(buff, v); err != nil {
		return logAndError(
			logger,
			fmt.Sprintf("response body: %s", string(buff)),
			fmt.Errorf(errInvalidJSONBody, url, err),
		)
	}
	return nil
}

func logAndError(logger *xlog.Logger, log string, e error) error {
	if log != "" {
		logger.Errorf("%s --- %s", e.Error(), log)
	} else {
		logger.Errorf(e.Error())
	}
	return e
}

func writeStructIntoMultipart(dst *multipart.Writer, v interface{}) error {
	val := reflect.ValueOf(v).Elem()
	typ := val.Type()
	numFields := val.NumField()
	for i := 0; i < numFields; i++ {
		var fieldName, fieldValue string
		var omitEmtpy bool
		fv := val.Field(i)
		tag := typ.Field(i).Tag.Get("json")
		if tag != "" {
			parts := xstrings.SplitAndTrim(tag, ",")
			fieldName = parts[0]
			if len(parts) > 1 {
				omitEmtpy = parts[1] == "omitempty"
			}
		}
		if !fv.IsValid() || (omitEmtpy && isEmptyValue(fv)) {
			continue
		}
		if fieldName == "" {
			fieldName = xstrings.ToSnakeCase(typ.Name())
		}
		iv := fv.Interface()
		switch iv.(type) {
		case int, int8, int16, int32, int64:
			if omitEmtpy && fv.Int() == 0 {
				continue
			}
			fieldValue = strconv.FormatInt(fv.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			if omitEmtpy && fv.Uint() == 0 {
				continue
			}
			fieldValue = strconv.FormatUint(fv.Uint(), 10)
		case []byte:
			fieldValue = fv.String()
		case string:
			fieldValue = fv.String()
		default:
			buff, err := json.Marshal(iv)
			if err != nil {
				return xerrors.Wrap(err, "could not convert %s to url.Values", typ.Name())
			}
			fieldValue = string(buff)
		}
		if err := dst.WriteField(fieldName, fieldValue); err != nil {
			return xerrors.Wrap(err, "could not write %q field (value=%s)", fieldName, fieldValue)
		}
	}
	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
