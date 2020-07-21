package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"gitlab.qn.com/ydbt/usetool"
)

// 配置热加载

/* ------------------- 配置文件 ------------------- */
// NewFileCfger
// cfgFile:配置文件；cfgType:配置文件类型；interval:热加载间隔
func NewFileCfger(cfgFile, cfgSuffix string, interval int) *FileCfger {
	fcr := &FileCfger{
		pollInterval: interval,
		cfgFile:      cfgFile,
		cfgSuffix:    cfgSuffix,
		cfgInfo:      NewProjectCfg(),
		subscribers:  make(map[string]distributeCfg),
	}
	if cfgSuffix == "" {
		index := strings.LastIndex(cfgFile, ".")
		if index == -1 {
			return nil
		} else {
			fcr.cfgSuffix = cfgFile[index+1:]
		}
	}
	fcr.LoadConfig()
	return fcr
}

// Regist
// 顺序初始化并未加锁
// 此处使用的具体实例类，并未使用接口
func (fcr *FileCfger) Regist(ti interface{}) {
	regist_subscribers(ti, fcr.subscribers)
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
			cfg, err := fcr.LoadConfig()
			if err != nil || cfg == nil {
				fmt.Println(err)
				continue
			}
			for _, v := range fcr.subscribers {
				v.funcSpecialCfg(cfg, v.hli)
			}
		}
	}
}

// LoadCfgFile
// 根据指定的文件类型加载配置文件
func (fcr *FileCfger) LoadConfig() (*ProjectCfg, error) {
	var err error
	cfgByte, err := ioutil.ReadFile(fcr.cfgFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	switch fcr.cfgSuffix {
	case "js", "json":
		err = Json2ProjectCfg(string(cfgByte), fcr.cfgInfo)
	case "yaml", "yml":
		err = Yaml2ProjectCfg(string(cfgByte), fcr.cfgInfo)
	}
	return fcr.cfgInfo, err
}

// Config
// 获取已经缓存的配置信息
func (fcr *FileCfger) Config() *ProjectCfg {
	return fcr.cfgInfo
}

// FileCfger
// 应用通过配置文件获取配置
type FileCfger struct {
	pollInterval int
	cfgFile      string
	cfgSuffix    string
	cfgInfo      *ProjectCfg
	subscribers  map[string]distributeCfg
}

type distributeCfg struct {
	hli            HotLoadI
	funcSpecialCfg func(cfg *ProjectCfg, hli HotLoadI)
}
