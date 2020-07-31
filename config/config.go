package config

func NewCfger(cfgType string, parsei ProjectConfigI, interval int) CfgerI {
	switch cfgType {
	case "cfgfile", "file":
		return NewFileCfger(parsei, interval) // 间隔10秒读取一下配置文件
	case "cfgserver", "cfgservice", "server", "service":
		return NewApolloCfger(parsei, interval)
	}
	return nil
}
