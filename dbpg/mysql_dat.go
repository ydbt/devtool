package dbpg

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// MysqlCfg
// MaxOpenConn 最大的数据库连接数
// MaxIdleConn 最大的空闲的数据库连接数
// MaxLifeTime 单位为秒
type MysqlCfg struct {
	Host        string `json:"host" yaml:"host"`
	Port        int    `json:"port" yaml:"port"`
	User        string `json:"user" yaml:"user"`
	Password    string `json:"password" yaml:"pasword"`
	Database    string `json:"database" yaml:"database"`
	Charset     string `json:"charset" yaml"charset"`
	MaxOpenConn int    `json:"maxopenconn" yaml:"maxopenconn"`
	MaxIdleConn int    `json:"maxidleconn" yaml:"maxidleconn"`
	MaxLifeTime int    `json:"maxlifetime" yaml:"maxlifetime"`
}

// newMysqlCfg
func newMysqlCfg() *MysqlCfg {
	return &MysqlCfg{
		Host:        "127.0.0.1",
		Port:        3306,
		User:        "sdadmin",
		Password:    "admin",
		Database:    "sdadmin",
		Charset:     "utf8",
		MaxOpenConn: 5,
		MaxIdleConn: 2,
		MaxLifeTime: 60 * 60,
	}
}

func (cfg *MysqlCfg) check() {
	if cfg.Port < 1024 || cfg.Port > 102400 {
		cfg.Port = 3306
	}
	if cfg.MaxOpenConn < 1 || cfg.MaxOpenConn > 100 {
		cfg.MaxOpenConn = 10
	}
	if cfg.MaxIdleConn < 1 || cfg.MaxIdleConn > 100 {
		cfg.MaxOpenConn = cfg.MaxOpenConn/4 + 1
	}
	if cfg.MaxLifeTime < 60 || cfg.MaxLifeTime > 2400*3600 {
		cfg.MaxLifeTime = 2400 * 3600
	}
}

func Json2MysqlCfg(js string) (*MysqlCfg, error) {
	cfg := newMysqlCfg()
	err := json.Unmarshal([]byte(js), cfg)
	cfg.check()
	return cfg, err
}

func Yaml2MysqlCfg(yml string) (*MysqlCfg, error) {
	cfg := newMysqlCfg()
	err := yaml.Unmarshal([]byte(yml), cfg)
	cfg.check()
	return cfg, err
}

type DbRow []string
type DbTable []DbRow
