package config

import "github.com/ydbt/devtool/usetool"

// HotLoadI
// 需要热加载配置的订阅制接口，转换为自己的配置对象
type HotLoadI interface {
	UpdateCfg(cfg interface{})
}

// SubscribeCfg
// 配置管理者，分发热加载配置
type SubscribeCfg interface {
	Regist(interface{}, *HotLoadI)
}

// LoadCfgI
// 配置管理者
type CfgerI interface {
	// 立即加载配置
	LoadConfig() (*ProjectCfg, error)
	// 获取缓存配置
	Config() *ProjectCfg
	// 定时获取配置信息
	TimerPollLoadCfg(psi usetool.ProcessSignalI)
}
