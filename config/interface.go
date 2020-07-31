package config

import "github.com/ydbt/devtool/v3/usetool"

// HotLoadI
// 需要热加载配置的订阅制接口，转换为自己的配置对象
type HotLoadI interface {
	UpdateCfg(cfg interface{})
}

type funcConfig func(interface{}) interface{}
type DynamicLoadCfg struct {
	// 配置接受者
	FuncHotLoad HotLoadI
	// 配置项
	GetCfg funcConfig
}

type ProjectConfigI interface {
	// 字节解码项目配置
	Marshal() (out []byte, err error)
	// 项目配置对象字节编码
	Unmarshal(data []byte) error
	// 获取项目配置
	Config() interface{}
	// 获取配置文件
	CfgFile() string
}

// LoadCfgI
// 配置管理者
type CfgerI interface {
	// 立即加载配置
	LoadConfig() error
	// 获取缓存配置，调用ProjectConfigI Config接口
	Config() interface{}
	// 定时获取配置信息
	TimerPollLoadCfg(psi usetool.ProcessSignalI)
	// 配置管理者，分发热加载配置，
	Regist(DynamicLoadCfg)
}
