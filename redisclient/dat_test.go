package redisclient_test

import (
	"testing"
	"gitlab.qn.com/ydbt/redisclient"
)

func TestRedisCfgInit(t *testing.T) {
	yml := `
addrs: ["127.0.0.1:6379","localhost:6379","0.0.0.0:6379"]
user: "ut-sdadmin"
password: "ut-admin"
database: -1
maxretry: -1
poolsize: 100
minidleconns: 10
pooltimeout: -1
idletimeout: -1
`
	ymlCfg, err := redisclient.Yaml2RedisCfg(yml)
	if err != nil {
		t.Error("redis config yaml parse failed: ", err)
		return
	}
	js := `
{
  "addrs": ["127.0.0.1:6379","localhost:6379","0.0.0.0:6379"],
  "user": "ut-sdadmin",
  "password": "ut-admin",
  "database": -1,
  "maxretry": -1,
  "poolsize": 100,
  "minidleconns": 10,
  "pooltimeout": -1,
  "idletimeout": -1
}
`
	jsCfg, err := redisclient.Json2RedisCfg(js)
	if err != nil {
		t.Error("redis config json parse failed: ", err)
		return
	}

	if (ymlCfg.User != jsCfg.User) && (jsCfg.User != "ut-sdadmin") {
		t.Error("parse config file failed")
		t.Error("yml: ", ymlCfg)
		t.Error("js: ", jsCfg)
	}

}
