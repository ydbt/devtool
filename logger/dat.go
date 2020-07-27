package logger

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
)

type FuncStrategy func() map[string]interface{}

// LogCfg 日志回滚文件配置信息
type LogCfg struct {
	AppName      string `json:"appname" yaml"appname"`
	Path         string `json:"path" yaml:"path"`
	Level        string `json:"level" yaml:"level"`
	MaxSize      int    `json:"maxsize" yaml:"maxsize"`
	MaxBackup    int    `json:"maxbackup" yaml:"maxbackup"`
	MaxLive      int    `json:"maxlive" yaml:"maxlive"`
	IsCompress   bool   `json:"compress" yaml:"compress"`
	Strategy     int    `json:"strategy" yaml:"strategy"`
	funcStrategy FuncStrategy
}

// trace(-1) < debug(0) < info(1) < warn(2) < error(3) < fatal(4) < panic(5)
// 默认等级为 info

func newLogCfg() *LogCfg {
	return &LogCfg{
		AppName:    "app",
		Path:       "app.log",
		Level:      "debug",
		MaxSize:    100,
		MaxBackup:  100,
		MaxLive:    60,
		IsCompress: false,
		Strategy:   0}
}

func Json2LogCfg(js string) (*LogCfg, error) {
	cfg := newLogCfg()
	err := json.Unmarshal([]byte(js), &cfg)
	if err != nil {
		return cfg, err
	}
	cfg.check()
	return cfg, nil
}

func Yaml2LogCfg(yml string) (*LogCfg, error) {
	cfg := newLogCfg()
	err := yaml.Unmarshal([]byte(yml), &cfg)
	if err != nil {
		return cfg, err
	}
	cfg.check()
	return cfg, nil
}

func (cfg *LogCfg) check() {
	if cfg.Strategy < 0 {
		cfg.Strategy = 0
	} else if cfg.Strategy > 3 {
		cfg.Strategy = 3
	}
	level := strings.ToLower(cfg.Level)
	if level != "trace" &&
		level != "debug" &&
		level != "info" &&
		level != "warn" &&
		level != "error" &&
		level != "fatal" {
		level = "debug"
	}
	cfg.Level = level
}

func strategy0() map[string]interface{} {
	return make(map[string]interface{})
}

func strategy1() map[string]interface{} {
	mapCallers := make(map[string]interface{})
	mapCallers["func"] = "#"
	pcs := make([]uintptr, 1)
	wsize := runtime.Callers(3, pcs)
	if wsize <= 0 {
		return mapCallers
	}
	pfs := runtime.CallersFrames(pcs)
	f, _ := pfs.Next()
	index := strings.LastIndex(f.Function, ".")
	if index == -1 {
		index = 0
	}
	mapCallers["func"] = f.Function[index+1:]
	return mapCallers
}

func strategy2() map[string]interface{} {
	mapCallers := make(map[string]interface{})
	mapCallers["func"] = "#"
	mapCallers["fileposition"] = "-:0"

	pcs := make([]uintptr, 1)
	wsize := runtime.Callers(3, pcs)
	if wsize <= 0 {
		return mapCallers
	}
	pfs := runtime.CallersFrames(pcs)
	f, _ := pfs.Next()
	index := strings.LastIndex(f.Function, ".")
	if index == -1 {
		index = 0
	}
	mapCallers["func"] = f.Function[index+1:]
	index = strings.LastIndex(f.File, "/")
	if index == -1 {
		index = 0
	}
	mapCallers["fileposition"] = fmt.Sprintf("%s:%d", f.File[index+1:], f.Line)
	return mapCallers
}

func strategy3() map[string]interface{} {
	mapCallers := make(map[string]interface{})
	mapCallers["func"] = "#"
	mapCallers["fileposition"] = "-:0"

	pcs := make([]uintptr, 1)
	wsize := runtime.Callers(3, pcs)
	if wsize <= 0 {
		return mapCallers
	}
	pfs := runtime.CallersFrames(pcs)
	f, _ := pfs.Next()
	mapCallers["func"] = f.Function
	mapCallers["fileposition"] = fmt.Sprintf("%s:%d", f.File, f.Line)
	return mapCallers
}
