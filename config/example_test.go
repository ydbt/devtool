package config_test

import (
	"fmt"
	"os"
	"time"

	"github.com/ydbt/devtool/v3/config"
	"github.com/ydbt/devtool/v3/logger"
	"github.com/ydbt/devtool/v3/usetool"
)

func Example_HotLoadFile() {
	cfg := &ymlProjectCfg{
		pro:  new(ProjectCfg),
		file: "ut_file_test.yml",
	}
	cfg.pro = new(ProjectCfg)
	psr := usetool.NewProcessSignaler()
	fcr := config.NewFileCfger(cfg, 1)
	err := fcr.LoadConfig()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	cfgInfo := fcr.Config()
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
		fmt.Printf("ProjectCfg fetch logger.Logcfg failed, %v\n", cfgInfo)
		return
	}
	lg := logger.NewLogger(cfgLog) // 日志热加载
	fcr.Regist(config.DynamicLoadCfg{
		FuncHotLoad: lg,
		GetCfg:      funcCfgLog,
	}) // 日志实例
	go fcr.TimerPollLoadCfg(psr)      // 热加载轮询
	lg.Debugf("********************") // 观测日志变化
	lg.Infof("###################")
	for i := 0; i < 10; i++ {
		lg.Infof("=======================")
		time.Sleep(time.Second * 2)
	}
	psr.WaitSignalProcess(os.Interrupt) // 进程终止，通知协程终止
}
