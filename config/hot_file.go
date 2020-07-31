package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ydbt/devtool/v3/usetool"
)

// 配置热加载

/* ------------------- 配置文件 ------------------- */
// NewFileCfger
// cfgFile:配置文件；cfgType:配置文件类型；interval:热加载间隔
func NewFileCfger(parsei ProjectConfigI, interval int) *FileCfger {
	fcr := &FileCfger{
		pollInterval: interval,
		cfgFile:      parsei.CfgFile(),
		parsei:       parsei,
	}
	fcr.LoadConfig()
	return fcr
}

// Regist
// 顺序初始化并未加锁
// 此处使用的具体实例类，并未使用接口
func (fcr *FileCfger) Regist(suber DynamicLoadCfg) {
	fcr.subscribers = append(fcr.subscribers, suber)
}

// TimerPollLoadCfg
// 定时轮询配置文件，热加载配置项
func (fcr *FileCfger) TimerPollLoadCfg(psi usetool.ProcessSignalI) {
	timerSignal := time.NewTimer(time.Second * time.Duration(fcr.pollInterval))
	interruptSignal := psi.TopicOSSignal(os.Interrupt)
	killSignal := psi.TopicOSSignal(os.Kill)
	// 缓存队列初始化
	for {
		fmt.Println("hot config poll")
		select {
		case <-interruptSignal:
			return
		case <-killSignal:
			return
		case <-timerSignal.C:
			timerSignal.Reset(time.Second * time.Duration(fcr.pollInterval))
			err := fcr.LoadConfig()
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, v := range fcr.subscribers {
				v.FuncHotLoad.UpdateCfg(v.GetCfg(fcr.parsei.Config()))
			}
		}
	}
}

// LoadCfgFile
// 根据指定的文件类型加载配置文件
func (fcr *FileCfger) LoadConfig() error {
	var err error
	cfgByte, err := ioutil.ReadFile(fcr.cfgFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = fcr.parsei.Unmarshal(cfgByte)
	return err
}

// Config
// 获取已经缓存的配置信息
func (fcr *FileCfger) Config() interface{} {
	return fcr.parsei.Config()
}

// FileCfger
// 应用通过配置文件获取配置
type FileCfger struct {
	pollInterval int
	cfgFile      string
	parsei       ProjectConfigI
	subscribers  []DynamicLoadCfg
}
