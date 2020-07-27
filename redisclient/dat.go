package redisclient

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// RedisCfg
// Addrs:["127.0.0.1:7000",...]
type RedisCfg struct {
	Addrs        []string `json:"addrs" yaml:"addrs"`
	User         string   `json:"user" yaml:"user"`
	Password     string   `json:"password" yaml:"password"`
	Database     int      `json:"database" yaml:"database"`
	MaxRetry     int      `json:"maxretry" yaml:"maxretry"`
	PoolSize     int      `json:"poolsize" yaml:"poolsize"`
	MinIdleConns int      `json:"minidleconns" yaml:"minidleconns"`
	PoolTimeout  int      `json:"pooltimeout" yaml:"pooltimeout"`
	IdleTimeout  int      `json:"idletimeout" yaml:"idletimeout"`
}

func newRedisCfg() *RedisCfg {
	return &RedisCfg{
		Addrs:        []string{"redis01:6379", "redis02:6379", "redis03:6379"},
		User:         "sdadmin",
		Password:     "",
		Database:     0,
		MaxRetry:     -1,
		PoolSize:     10,
		MinIdleConns: 5,
		PoolTimeout:  60000,
		IdleTimeout:  60000,
	}
}

func Json2RedisCfg(js string) (*RedisCfg, error) {
	cfg := newRedisCfg()
	err := json.Unmarshal([]byte(js), &cfg)
	return cfg, err
}

func Yaml2RedisCfg(yml string) (*RedisCfg, error) {
	cfg := newRedisCfg()
	err := yaml.Unmarshal([]byte(yml), &cfg)
	return cfg, err
}
