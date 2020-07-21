package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ydbt/devtool/usetool"
	"gopkg.in/yaml.v2"

	agolloV "github.com/zouyx/agollo/v3"
	"github.com/zouyx/agollo/v3/env/config"
	agolloVConfig "github.com/zouyx/agollo/v3/env/config"
)

/* ------------------- Apollo 配置 ------------------- */ // TODO
func NewApolloCfger(cfgFile, cfgSuffix string, interval int) *ApolloCfger {
	acr := &ApolloCfger{
		pollInterval: interval,
		cfgFile:      cfgFile,
		cfgSuffix:    cfgSuffix,
		cfgInfo:      NewProjectCfg(),
		cfgApollo:    new(ApolloCfg),
		subscribers:  make(map[string]distributeCfg),
	}
	if cfgSuffix == "" {
		index := strings.LastIndex(cfgFile, ".")
		if index == -1 {
			return nil
		} else {
			acr.cfgSuffix = cfgFile[index+1:]
		}
	}
	acr.LoadConfig()
	return acr
}

// LoadCfgFile
// 根据指定的文件类型加载配置文件
func (acr *ApolloCfger) LoadConfig() (*ProjectCfg, error) {
	var err error
	cfgByte, err := ioutil.ReadFile(acr.cfgFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	switch acr.cfgSuffix {
	case "js", "json":
		err := json.Unmarshal([]byte(cfgByte), acr.cfgApollo)
		if err != nil {
			return nil, err
		}
	case "yaml", "yml":
		err := yaml.Unmarshal([]byte(cfgByte), acr.cfgApollo)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("read file:", acr.cfgApollo)
	ymlCfg := acr.pullConfig()
	Yaml2ProjectCfg(ymlCfg, acr.cfgInfo)
	return acr.cfgInfo, nil
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
			_, err := acr.LoadConfig()
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, v := range acr.subscribers {
				v.funcSpecialCfg(acr.cfgInfo, v.hli)
			}
		}
	}
}

func (acr *ApolloCfger) pullConfig() string {
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
	fmt.Println(ymlCfg)
	fmt.Println("pullConfig end")
	return ymlCfg
}

func (acr *ApolloCfger) Config() *ProjectCfg {
	return acr.cfgInfo
}

// Regist
// 顺序初始化并未加锁
// 此处使用的具体实例类，并未使用接口
func (acr *ApolloCfger) Regist(ti interface{}) {
	regist_subscribers(ti, acr.subscribers)
}

type ApolloCfger struct {
	pollInterval int
	cfgFile      string
	cfgSuffix    string
	cfgApollo    *ApolloCfg
	cfgInfo      *ProjectCfg
	subscribers  map[string]distributeCfg
}
