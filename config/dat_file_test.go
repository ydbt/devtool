package config_test

import (
	"testing"

	"gitlab.qn.com/ydbt/config"
)

func TestProjectCfgYaml(t *testing.T) {
	yml := `
services:
  -
    host: 192.168.1.1
    port: 100861
  -
    host: 127.0.0.1
    port: 10010
log:
  appname: config-log
  path: config.log
  level: trace
  maxsize: 10
  maxbackup: 10
  maxlive: 60
  compress: true
  strategy: 2
mysql:
  host: "127.0.0.1"
  port: 3306
  user: "config-sdadmin"
  password: "config-admin"
  database: "config-sdadmin"
  charset: "utf8"
  maxopenconn: 5
  maxidleconn: 2
  maxlifetime: 3600
redis:
  addrs: ["127.0.0.1:6379","localhost:6379","0.0.0.0:6379"]
  user: "config-sdadmin"
  password: "config-admin"
  database: -1
  maxretry: -1
  poolsize: 100
  minidleconns: 10
  pooltimeout: -1
  idletimeout: -1
`
	cfg := config.NewProjectCfg()
	err := config.Yaml2ProjectCfg(yml, cfg)
	if err != nil {
		t.Error(err)
		return
	}
	if len(cfg.Redis.Addrs) != 3 {
		t.Error("redis addrs list parse failed")
		return
	}
	if cfg.Redis.User != "config-sdadmin" {
		t.Errorf("\"config-sdadmin\" != \"%s\"", cfg.Redis.User)
		return
	}
}

func TestProjectCfgJson(t *testing.T) {
	js := `
{
  "services":[{
      "host":"192.168.1.1",
      "port":1008611
    },
    {
      "host":"127.0.0.1",
      "port":10010
    }],
  "log":
  {
    "appname": "config-log",
    "path": "ut.log",
    "level": "trace",
    "maxsize": 100,
    "maxbackup": 200,
    "maxlive": 60,
    "compress": true,
    "strategy": 2
  },
  "mysql":
  {
    "host": "127.0.0.1",
    "port": 3306,
    "user": "config-sdadmin",
    "password": "config-admin",
    "database": "config-sdadmin",
    "charset": "utf8",
    "maxopenconn": 5,
    "maxidleconn": 2,
    "maxlifetime": 3600  
  },
  "redis":
  {
    "addrs": ["192.168.1.1:6379","127.0.0.1:6379"],
    "user": "config-sdadmin",
    "password": "config-admin",
    "database": -1,
    "maxretry": -1,
    "poolsize": 100,
    "minidleconns": 10,
    "pooltimeout": -1,
    "idletimeout": -1
  }
}
`
	cfg := config.NewProjectCfg()
	err := config.Json2ProjectCfg(js, cfg)
	if err != nil {
		t.Error(err)
		return
	}
	if len(cfg.Redis.Addrs) != 2 {
		t.Error("redis addrs list parse failed")
		return
	}
	if cfg.Redis.User != "config-sdadmin" {
		t.Errorf("\"config-sdadmin\" != \"%s\"", cfg.Redis.User)
		return
	}
}
