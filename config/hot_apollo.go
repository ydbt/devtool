package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ydbt/devtool/v3/usetool"
	"gopkg.in/yaml.v2"

	agolloV "github.com/zouyx/agollo/v3"
	"github.com/zouyx/agollo/v3/env/config"
	agolloVConfig "github.com/zouyx/agollo/v3/env/config"
)

/* ------------------- Apollo 配置 ------------------- */ // TODO
func NewApolloCfger(parsei ProjectConfigI, interval int) *ApolloCfger {
	acr := &ApolloCfger{
		pollInterval: interval,
		cfgFile:      parsei.CfgFile(),
		cfgApollo:    new(ApolloCfg),
		parsei:       parsei,
	}
	index := strings.LastIndex(acr.cfgFile, ".")
	if index == -1 {
		return nil
	} else {
		acr.cfgSuffix = acr.cfgFile[index+1:]
	}
	acr.LoadConfig()
	return acr
}

// LoadCfgFile
// 根据指定的文件类型加载配置文件
func (acr *ApolloCfger) LoadConfig() error {
	var err error
	cfgByte, err := ioutil.ReadFile(acr.cfgFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	switch acr.cfgSuffix {
	case "js", "json":
		err := json.Unmarshal([]byte(cfgByte), acr.cfgApollo)
		if err != nil {
			return err
		}
	case "yaml", "yml":
		err := yaml.Unmarshal([]byte(cfgByte), acr.cfgApollo)
		if err != nil {
			return err
		}
	}
	fmt.Println("read file:", acr.cfgApollo)
	ymlCfg := acr.pullConfig()
	acr.parsei.Unmarshal(ymlCfg)
	return nil
}

func (acr *ApolloCfger) GoTimerPollLoadCfg(psi usetool.ProcessSignalI) {
	go acr.TimerPollLoadCfg(psi)
}

func (acr *ApolloCfger) TimerPollLoadCfg(psi usetool.ProcessSignalI) {
	timerSignal := time.NewTimer(time.Second * time.Duration(acr.pollInterval))
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
			timerSignal.Reset(time.Second * time.Duration(acr.pollInterval))
			err := acr.LoadConfig()
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, v := range acr.subscribers {
				v.FuncHotLoad.UpdateCfg(v.GetCfg(acr.parsei.Config()))
			}
		}
	}
}

func (acr *ApolloCfger) pullConfig() []byte {
	fmt.Println("pullConfig beg")
	mapCfg := make(map[string]string)
	for _, spcName := range acr.cfgApollo.Namespaces {
		c := &agolloVConfig.AppConfig{
			AppID:          acr.cfgApollo.AppID,
			Cluster:        acr.cfgApollo.Cluster,
			IP:             acr.cfgApollo.IP,
			NamespaceName:  spcName,
			IsBackupConfig: false,
			Secret:         acr.cfgApollo.Secret,
		}
		agolloV.InitCustomConfig(func() (*config.AppConfig, error) {
			return c, nil
		})
		agolloV.Start()
		cache := agolloV.GetConfigCache(spcName)
		cache.Range(func(key, value interface{}) bool {
			mapCfg[fmt.Sprint(key)] = fmt.Sprint(value)
			return true
		})
	}
	ymlCfg := Apollo2Yaml(mapCfg)
	//	fmt.Println(ymlCfg)
	fmt.Println("pullConfig end")
	return []byte(ymlCfg)
}

func (acr *ApolloCfger) Config() interface{} {
	return acr.parsei.Config()
}

// Regist
// 顺序初始化并未加锁
// 此处使用的具体实例类，并未使用接口
func (acr *ApolloCfger) Regist(suber DynamicLoadCfg) {
	acr.subscribers = append(acr.subscribers, suber)
}

type ApolloCfger struct {
	pollInterval int
	cfgFile      string
	cfgSuffix    string
	cfgApollo    *ApolloCfg
	parsei       ProjectConfigI
	subscribers  []DynamicLoadCfg
}
