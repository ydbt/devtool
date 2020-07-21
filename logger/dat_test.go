package logger_test

import (
	"testing"
	"github.com/ydbt/devtool/v3/logger"
)

func TestLogCfgInit(t *testing.T) {
	yml := `
appname: ut-log
path: ut.log
level: trace
maxsize: 10
maxbackup: 10
maxlive: 60
compress: true
strategy: 2
`
	ymlCfg, err := logger.Yaml2LogCfg(yml)
	if err != nil {
		t.Error("log config yaml parse failed, ", err)
		return
	}
	js := `
{
  "appname": "ut-log",
  "path": "ut.log",
  "level": "trace",
  "maxsize": 100,
  "maxbackup": 200,
  "maxlive": 60,
  "compress": true,
  "strategy": 2
}
`
	jsCfg, err := logger.Json2LogCfg(js)
	if err != nil {
		t.Error("log config json parse failed, ", err)
		return
	}
	if (jsCfg.AppName != ymlCfg.AppName) && (jsCfg.AppName != "ut-log") {
		t.Error(jsCfg, ymlCfg)
	}
	if (jsCfg.Strategy != ymlCfg.Strategy) && (jsCfg.Strategy != 2) {
		t.Error(jsCfg, ymlCfg)
	}

}
