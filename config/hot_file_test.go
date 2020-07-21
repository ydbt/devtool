package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/ydbt/devtool/config"
	"github.com/ydbt/devtool/logger"
	"github.com/ydbt/devtool/usetool"
)

// TestHotLoadJson
// 启动测试通过修改tmp_ut.yml观察日志输出界别或则数据策略
func TestHotFileLoadJson(t *testing.T) {
	psr := usetool.NewProcessSignaler()
	fcr := config.NewFileCfger("ut_file_test.yml", "yml", 1)
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
