package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/ydbt/devtool/v3/config"
	"github.com/ydbt/devtool/v3/logger"
	"github.com/ydbt/devtool/v3/usetool"
)

// TestHotLoadJson
// 启动测试通过修改tmp_ut.yml观察日志输出界别或则数据策略
func TestHotApolloLoadYaml(t *testing.T) {
	cfg := &ymlProjectCfg{
		pro:  new(ProjectCfg),
		file: "ut_apollo_file.yml",
	}
	cfg.pro = new(ProjectCfg)

	psr := usetool.NewProcessSignaler()
	acr := config.NewApolloCfger(cfg, 5)
	err := acr.LoadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	cfgInfo := acr.Config()
	funcCfgLog := func(cfgPro interface{}) interface{} {
		var cfgLog *logger.LogCfg
		switch cfgPro := cfgInfo.(type) {
		case ProjectCfg:
			cfgLog = &cfgPro.Log
		case *ProjectCfg:
			cfgLog = &cfgPro.Log
		default:
			cfgLog = nil
		}
		return cfgLog
	}
	var cfgLog *logger.LogCfg
	var ok bool
	if cfgLog, ok = funcCfgLog(cfgInfo).(*logger.LogCfg); !ok {
		t.Errorf("ProjectCfg fetch logger.Logcfg failed, %v", cfgInfo)
		return
	}
	lg := logger.NewLogger(cfgLog) // 日志热加载
	acr.Regist(config.DynamicLoadCfg{
		FuncHotLoad: lg,
		GetCfg:      funcCfgLog,
	}) // 日志实例
	go acr.TimerPollLoadCfg(psr)      // 热加载轮询
	lg.Debugf("********************") // 观测日志变化
	lg.Infof("###################")
	for i := 0; i < 30; i++ {
		lg.Infof("=======================")
		time.Sleep(time.Second * 2)
	}
	psr.WaitSignalProcess(os.Interrupt) // 进程终止，通知协程终止
}
