package config

import (
	"os"
	"testing"

	"google.golang.org/appengine"

	"github.com/speedland/go/gae/gaetest"
	"github.com/speedland/go/x/xtesting/assert"
)

func TestMain(m *testing.M) {
	Global("urlfetch_deadline", "30", "urlfetch deadline secondsa")
	Global("facebook_app_id", "", "Facebook App ID")
	os.Exit(gaetest.Run(func() int {
		return m.Run()
	}))
}

func TestConfig_All(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.CleanupStorage(gaetest.NewContext(), "", "myapp"))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestConfig.json", nil))
	appCtx, _ := appengine.Namespace(gaetest.NewContext(), "myapp")
	c := New()
	c.Register("myconfig", "mydefaultvalue", "app custom config")
	c.Register("myconfig2", "mydefaultvalue2", "app custom config2")

	configs := c.All(appCtx)
	a.EqInt(len(c.defaultKeys), len(configs))

	configMap := make(map[string]*ServiceConfig)
	for _, cfg := range configs {
		configMap[cfg.Key] = cfg
	}
	a.EqStr("45", configMap["urlfetch_deadline"].Value)
	a.EqStr("local-app-id", configMap["facebook_app_id"].Value)
	a.EqStr("datastore value", configMap["myconfig"].Value)
	a.EqStr("mydefaultvalue2", configMap["myconfig2"].Value)
}

func TestConfig_Get(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.CleanupStorage(gaetest.NewContext(), "", "myapp"))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestConfig.json", nil))
	appCtx, _ := appengine.Namespace(gaetest.NewContext(), "myapp")
	c := New()
	c.Register("myconfig", "mydefaultvalue", "app custom config")
	c.Register("myconfig2", "mydefaultvalue2", "app custom config2")

	a.EqStr("45", c.Get(appCtx, "urlfetch_deadline").Value)
	a.EqStr("local-app-id", c.Get(appCtx, "facebook_app_id").Value)
	a.EqStr("datastore value", c.Get(appCtx, "myconfig").Value)
	a.EqStr("mydefaultvalue2", c.Get(appCtx, "myconfig2").Value)
}

func TestConfig_Get_Fallback(t *testing.T) {
	a := assert.New(t)
	a.Nil(gaetest.CleanupStorage(gaetest.NewContext(), "", "myapp"))
	a.Nil(gaetest.FixtureFromFile(gaetest.NewContext(), "./fixtures/TestConfig.json", nil))
	appCtx, _ := appengine.Namespace(gaetest.NewContext(), "myapp")
	c := New()
	cfg := c.Get(appCtx, "urlfetch_deadline")
	cfg.Value = "50"
	c.Set(appCtx, cfg)
	a.EqStr("50", c.Get(appCtx, "urlfetch_deadline").Value)
	a.EqStr("45", c.Get(appCtx, "urlfetch_deadline").GlobalValue)
	a.EqStr("30", c.Get(appCtx, "urlfetch_deadline").DefaultValue)
}