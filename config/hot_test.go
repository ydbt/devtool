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
func TestHotLoadJson(t *testing.T) {
	psr := usetool.NewProcessSignaler()
	fcr := config.NewFileCfger("tmp_ut.yml", "yml", 1)
	cfg, err := fcr.LoadConfig()
	if err != nil || cfg == nil {
		t.Error(err)
		return
	}
	lg := logger.NewLogger(&cfg.Log)
	fcr.Regist(lg)
	go fcr.TimerPollLoadCfg(psr)
	lg.Debugf("********************")
	lg.Infof("###################")
	for i := 0; i < 20; i++ {
		lg.Infof("=======================")
		time.Sleep(time.Second * 2)
	}
	psr.WaitSignalProcess(os.Interrupt)
}
