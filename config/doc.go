package config

// json 配置文件
/**
{
    "services": [{
        "host": "192.168.1.1",
        "port": 1008611
    },{
        "host": "127.0.0.1",
        "port": 10010
    }],
    "log": {
        "appname": "config-log",
        "path": "ut.log",
        "level": "trace",
        "maxsize": 100,
        "maxbackup": 200,
        "maxlive": 60,
        "compress": true,
        "strategy": 2
    },
    "mysql": {
        "host": "127.0.0.1",
        "port": 3306,
        "user": "config-sdadmin",
        "password": "config-admin",
        "database": "config-sdadmin",
        "charset": "utf8",
        "maxopenconn": 5,
        "maxidleconn": 2,
        "maxlifetime": 3600
    },
    "redis": {
        "addrs": ["192.168.1.1:6379", "127.0.0.1:6379"],
        "user": "config-sdadmin",
        "password": "config-admin",
        "database": -1,
        "maxretry": -1,
        "poolsize": 100,
        "minidleconns": 10,
        "pooltimeout": -1,
        "idletimeout": -1
    }
}
*/

// yaml 配置文件
/**
services:
  -
    host: 192.168.1.1
    port: 100861
  -
    host: 127.0.0.1
    port: 10010
log:
  appname: config-log
  path: config.log
  level: trace
  maxsize: 10
  maxbackup: 10
  maxlive: 60
  compress: true
  strategy: 1
mysql:
  host: "127.0.0.1"
  port: 3306
  user: "config-sdadmin"
  password: "config-admin"
  database: "config-sdadmin"
  charset: "utf8"
  maxopenconn: 5
  maxidleconn: 2
  maxlifetime: 3600
redis:
  addrs: ["127.0.0.1:6379","localhost:6379","0.0.0.0:6379"]
  user: "config-sdadmin"
  password: "config-admin"
  database: -1
  maxretry: -1
  poolsize: 100
  minidleconns: 10
  pooltimeout: -1
  idletimeout: -1
*/

/** apollo
services.[0].host=192.168.1.1
services.[0].port=100861
services.[1].host=127.0.0.1
services.[1].port=10010
log.appname=config-log
log.path=config.log
log.level= trace
log.maxsize=10
log.maxbackup=10
log.maxlive=60
log.compress=true
log.strategy=1
mysql.host=127.0.0.1
mysql.port=3306
mysql.user=config-sdadmin
mysql.password=config-admin
mysql.database=config-sdadmin
mysql.charset=utf8
mysql.maxopenconn=5
mysql.maxidleconn=2
mysql.maxlifetime=3600
redis.addrs.[0]=127.0.0.1:6379
redis.addrs.[1]=localhost:6379
redis.addrs.[2]=0.0.0.0:637
redis.user=config-sdadmin
redis.password=config-admin
redis.database=-1
redis.maxretry=-1
redis.poolsize=100
redis.minidleconns=10
redis.pooltimeout=-1
redis.idletimeout=-1
*/
