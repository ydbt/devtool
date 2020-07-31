package config_test

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

type ymlProjectCfg struct {
	pro  *ProjectCfg
	file string
}

func (cfg *ymlProjectCfg) CfgFile() string {
	return cfg.file
}

func (cfg *ymlProjectCfg) Marshal() (out []byte, err error) {
	//	fmt.Println("yaml Marshal:\n", string(out))
	byteYml, err := yaml.Marshal(cfg.pro)
	return byteYml, err
}

func (cfg *ymlProjectCfg) Unmarshal(data []byte) error {
	//	fmt.Println("yaml Unmarshal:\n", string(data))
	err := yaml.Unmarshal(data, cfg.pro)
	return err
}

func (cfg *ymlProjectCfg) Config() interface{} {
	return cfg.pro
}

type jsonProjectCfg struct {
	pro  *ProjectCfg
	file string
}

func (cfg *jsonProjectCfg) Marshal() (out []byte, err error) {
	byteYml, err := json.Marshal(cfg.pro)
	return byteYml, err
}

func (cfg *jsonProjectCfg) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, cfg.pro)
	return err
}

func (cfg *jsonProjectCfg) Config() interface{} {
	return cfg.pro
}

func (cfg *jsonProjectCfg) CfgFile() string {
	return cfg.file
}
