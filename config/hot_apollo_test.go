package config_test

import (
	"os"
	"testing"
	"time"

	"gitlab.qn.com/ydbt/config"
	"gitlab.qn.com/ydbt/logger"
	"gitlab.qn.com/ydbt/usetool"
)

// TestHotLoadJson
// 启动测试通过修改tmp_ut.yml观察日志输出界别或则数据策略
func TestHotApolloLoadJson(t *testing.T) {
	psr := usetool.NewProcessSignaler()
	acr := config.NewApolloCfger("tmp_apollo_ut.yml", "yml", 10)
	cfg, err := acr.LoadConfig()
	if err != nil || cfg == nil {
		t.Error(err)
		return
	}
	lg := logger.NewLogger(&cfg.Log)
	acr.Regist(lg)
	go acr.TimerPollLoadCfg(psr)
	lg.Debugf("********************")
	lg.Infof("###################")
	for i := 0; i < 20; i++ {
		lg.Infof("=======================")
		time.Sleep(time.Second * 2)
	}
	psr.WaitSignalProcess(os.Interrupt)
}
