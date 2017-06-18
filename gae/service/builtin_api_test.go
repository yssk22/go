package service

import (
	"net/url"
	"testing"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/gae/service/config"
	"github.com/speedland/go/web/httptest"
	"github.com/speedland/go/web/response"
)

func newTestService() *Service {
	s := New("myapp")
	s.Config.Register("myconfig", "myconfigvalue", "custom config")
	DefaultBuiltinAPIConfig.ActivateEndpoints(s)
	return s
}

func Test_API_Configs_List(t *testing.T) {
	a := httptest.NewAssert(t)
	s := newTestService()
	a.Nil(gaetest.CleanupStorage(gaetest.NewContext(), "", s.Namespace()))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/Test_API_Configs.json", nil))

	recorder := gaetest.NewRecorder(s)
	var configs []*config.ServiceConfig
	res := recorder.TestGet("/myapp/admin/api/configs/")
	a.Status(response.HTTPStatusOK, res)
	a.JSON(&configs, res)
	a.OK(len(configs) > 0)
}

func Test_API_Configs_Get(t *testing.T) {
	a := httptest.NewAssert(t)
	s := newTestService()
	a.Nil(gaetest.CleanupStorage(gaetest.NewContext(), "", s.Namespace()))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/Test_API_Configs.json", nil))

	recorder := gaetest.NewRecorder(s)

	var cfg config.ServiceConfig
	res := recorder.TestGet("/myapp/admin/api/configs/urlfetch.deadline.json")
	a.Status(response.HTTPStatusOK, res)
	a.JSON(&cfg, res)
	a.EqStr("45", cfg.Value)

	res = recorder.TestGet("/myapp/admin/api/configs/myconfig.json")
	a.Status(response.HTTPStatusOK, res)
	a.JSON(&cfg, res)
	a.EqStr("datastore value", cfg.Value)
}

func Test_API_Configs_Put(t *testing.T) {
	a := httptest.NewAssert(t)
	s := newTestService()
	a.Nil(gaetest.CleanupStorage(gaetest.NewContext(), "", s.Namespace()))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/Test_API_Configs.json", nil))

	recorder := gaetest.NewRecorder(s)

	var cfg config.ServiceConfig
	res := recorder.TestPut(
		"/myapp/admin/api/configs/urlfetch.deadline.json",
		url.Values{
			"value": []string{"20"},
		},
	)
	a.Status(response.HTTPStatusOK, res)
	a.JSON(&cfg, res)
	a.EqStr("20", cfg.Value)

	res = recorder.TestGet("/myapp/admin/api/configs/urlfetch.deadline.json")
	a.Status(response.HTTPStatusOK, res)
	a.JSON(&cfg, res)
	a.EqStr("20", cfg.Value)

	// Check the global config value is not changed.
	globalCfg := config.DefaultServiceConfigKind.MustGet(gaetest.NewContext(), "urlfetch.deadline")
	a.NotNil(globalCfg)
	a.EqStr("45", globalCfg.Value)
}