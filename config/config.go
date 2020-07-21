package config

import (
	"fmt"

	"github.com/ydbt/devtool/v3/dbpg"
	"github.com/ydbt/devtool/v3/logger"
	"github.com/ydbt/devtool/v3/redisclient"
)

func NewCfger(cfgFile, cfgType, cfgSuffix string, interval int) CfgerI {
	switch cfgType {
	case "cfgfile", "file":
		return NewFileCfger(cfgFile, cfgSuffix, 10) // 间隔10秒读取一下配置文件
	case "cfgserver", "cfgservice", "server", "service":
		return NewApolloCfger(cfgFile, cfgSuffix, 10) // 间隔10秒读取一下配置文件
	}
	return nil
}

// Regist
// 顺序初始化并未加锁
// 此处使用的具体实例类，并未使用接口
func regist_subscribers(ti interface{}, subscribers map[string]distributeCfg) {
	switch ti.(type) {
	case logger.Logger:
		lg := ti.(logger.Logger)
		fmt.Println("HotLoad Regist logger.Logger")
		subscribers["log"] = distributeCfg{
			hli: &lg,
			funcSpecialCfg: func(cfg *ProjectCfg, hli HotLoadI) {
				hli.UpdateCfg(cfg.Log)
			},
		}
	case *logger.Logger:
		plg := ti.(*logger.Logger)
		fmt.Println("HotLoad Regist logger.Logger")
		subscribers["log"] = distributeCfg{
			hli: plg,
			funcSpecialCfg: func(cfg *ProjectCfg, hli HotLoadI) {
				hli.UpdateCfg(cfg.Log)
			},
		}
	case redisclient.RedisClient:
		rc := ti.(redisclient.RedisClient)
		subscribers["redis"] = distributeCfg{
			hli: &rc,
			funcSpecialCfg: func(cfg *ProjectCfg, hli HotLoadI) {
				hli.UpdateCfg(cfg.Redis)
			},
		}
	case *redisclient.RedisClient:
		prc := ti.(*redisclient.RedisClient)
		subscribers["redis"] = distributeCfg{
			hli: prc,
			funcSpecialCfg: func(cfg *ProjectCfg, hli HotLoadI) {
				hli.UpdateCfg(cfg.Redis)
			},
		}
	case dbpg.DbMysql:
		db := ti.(dbpg.DbMysql)
		subscribers["dbpg"] = distributeCfg{
			hli: &db,
			funcSpecialCfg: func(cfg *ProjectCfg, hli HotLoadI) {
				hli.UpdateCfg(cfg.Mysql)
			},
		}
	case *dbpg.DbMysql:
		pdb := ti.(*dbpg.DbMysql)
		subscribers["dbpg"] = distributeCfg{
			hli: pdb,
			funcSpecialCfg: func(cfg *ProjectCfg, hli HotLoadI) {
				hli.UpdateCfg(cfg.Mysql)
			},
		}
	}
}
