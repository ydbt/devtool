package config_test

import (
	"os"
	"time"

	"github.com/ydbt/devtool/v3/config"
	"github.com/ydbt/devtool/v3/logger"
	"github.com/ydbt/devtool/v3/usetool"
)

func Example_HotLoadFile() {
	psr := usetool.NewProcessSignaler()                // 进程信号管理
	fcr := config.NewFileCfger("tmp_ut.yml", "yml", 1) // 配置加载方式
	cfg, err := fcr.LoadConfig()                       // 加载配置文件
	if err != nil || cfg == nil {
		return
	}
	lg := logger.NewLogger(&cfg.Log)  // 日志热加载
	fcr.Regist(lg)                    // 日志实例
	go fcr.TimerPollLoadCfg(psr)      // 热加载轮询
	lg.Debugf("********************") // 观测日志变化
	lg.Infof("###################")
	for i := 0; i < 20; i++ {
		lg.Infof("=======================")
		time.Sleep(time.Second * 2)
	}
	psr.WaitSignalProcess(os.Interrupt) // 进程终止，通知协程终止
}
