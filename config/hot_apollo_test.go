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
func TestHotApolloLoadJson(t *testing.T) {
	psr := usetool.NewProcessSignaler()
	acr := config.NewApolloCfger("ut_apollo_file.yml", "yml", 10)
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
