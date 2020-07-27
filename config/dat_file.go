package config

import (
	"encoding/json"

	"github.com/ydbt/devtool/v3/dbpg"
	"github.com/ydbt/devtool/v3/kafkaclient"
	"github.com/ydbt/devtool/v3/logger"
	"github.com/ydbt/devtool/v3/redisclient"

	"gopkg.in/yaml.v2"
)

type ServiceCfg struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

type ProjectCfg struct {
	Services      []ServiceCfg                    `json:"services" yaml:"services"`
	Mysql         dbpg.MysqlCfg                   `json:"mysql" yaml:"mysql"`
	Redis         redisclient.RedisCfg            `json:"redis" yaml:"redis"`
	Log           logger.LogCfg                   `json:"log" yaml:"log"`
	KafkaProducer kafkaclient.KafkaProducerConfig `json:"kafkaproducer" yaml:"kafkaproducer"`
}

func NewProjectCfg() *ProjectCfg {
	return &ProjectCfg{
		Services: []ServiceCfg{ServiceCfg{
			Host: "127.0.0.1",
			Port: 10086,
		}},
		Mysql: dbpg.MysqlCfg{
			Host:        "127.0.0.1",
			Port:        3306,
			User:        "sdadmin",
			Password:    "admin",
			Database:    "sdadmin",
			Charset:     "utf8",
			MaxOpenConn: 5,
			MaxIdleConn: 2,
			MaxLifeTime: 60 * 60,
		},
		Redis: redisclient.RedisCfg{
			Addrs:        []string{"redis01:6379", "redis02:6379", "redis03:6379"},
			User:         "sdadmin",
			Password:     "",
			Database:     0,
			MaxRetry:     -1,
			PoolSize:     10,
			MinIdleConns: 5,
			PoolTimeout:  60000,
			IdleTimeout:  60000,
		},
		Log: logger.LogCfg{
			AppName:    "app",
			Path:       "app.log",
			Level:      "debug",
			MaxSize:    100,
			MaxBackup:  100,
			MaxLive:    60,
			IsCompress: false,
			Strategy:   0,
		},
		KafkaProducer: kafkaclient.KafkaProducerConfig{
			ServiceConfig: kafkaclient.KafkaConfig{
				BrokerList: []string{"127.0.0.1:9092", "127.0.0.1:9093", "127.0.0.1:9094"},
				ApiVersion: "2.5.0",
			},
			CustomConfig: kafkaclient.KafkaProducerCustomConfig{
				BufferSize: 128,
			},
		},
	}
}

func Json2ProjectCfg(js string, cfg *ProjectCfg) error {
	err := json.Unmarshal([]byte(js), cfg)
	if err != nil {
		return err
	}
	jsServices, _ := json.Marshal(cfg.Services)
	jsLog, _ := json.Marshal(cfg.Log)
	jsMysql, _ := json.Marshal(cfg.Mysql)
	jsRedis, _ := json.Marshal(cfg.Redis)
	jsKafka, _ := json.Marshal(cfg.KafkaProducer)
	CreateProjectCfgByJson(string(jsServices), string(jsLog), string(jsMysql), string(jsRedis), string(jsKafka), cfg)
	return err
}

func Yaml2ProjectCfg(yml string, cfg *ProjectCfg) error {
	err := yaml.Unmarshal([]byte(yml), cfg)
	if err != nil {
		return err
	}
	ymlServices, _ := yaml.Marshal(cfg.Services)
	ymlLog, _ := yaml.Marshal(cfg.Log)
	ymlMysql, _ := yaml.Marshal(cfg.Mysql)
	ymlRedis, _ := yaml.Marshal(cfg.Redis)
	ymlKafka, _ := yaml.Marshal(cfg.KafkaProducer)
	CreateProjectCfgByYaml(string(ymlServices), string(ymlLog), string(ymlMysql), string(ymlRedis), string(ymlKafka), cfg)
	return err
}

func CreateProjectCfgByJson(jsService, jsLog, jsMysql, jsRedis, jsKfkProducer string, cfg *ProjectCfg) {
	json.Unmarshal([]byte(jsService), &cfg.Services)
	cfgLog, _ := logger.Json2LogCfg(jsLog)
	cfgMysql, _ := dbpg.Json2MysqlCfg(jsMysql)
	cfgRedis, _ := redisclient.Json2RedisCfg(jsRedis)
	cfgKfkProducer, _ := kafkaclient.Json2KafkaProducerCfg(jsKfkProducer)
	cfg.Log = *cfgLog
	cfg.Mysql = *cfgMysql
	cfg.Redis = *cfgRedis
	cfg.KafkaProducer = *cfgKfkProducer
}

func CreateProjectCfgByYaml(ymlService, ymlLog, ymlMysql, ymlRedis, ymlKfkProducer string, cfg *ProjectCfg) {
	yaml.Unmarshal([]byte(ymlService), &cfg.Services)
	cfgLog, _ := logger.Yaml2LogCfg(ymlLog)
	cfgMysql, _ := dbpg.Yaml2MysqlCfg(ymlMysql)
	cfgRedis, _ := redisclient.Yaml2RedisCfg(ymlRedis)
	cfgKfkProducer, _ := kafkaclient.Yaml2KafkaProducerCfg(ymlKfkProducer)
	cfg.Log = *cfgLog
	cfg.Mysql = *cfgMysql
	cfg.Redis = *cfgRedis
	cfg.KafkaProducer = *cfgKfkProducer
}
