# go 常用第三方包接口封装
require github.com/ydbt/devtool/v3/v1 v1.0.0
## 测试方案
### 持续测试
### 压力测试
### 功能指标
## 编码规则
### 命名规则
### 注释文档
## 工具集
### 日志
1. pkg/logger 为日志封装实现
2. 使用[zerolog](https://github.com/rs/zerolog.git)作为底层实现层
### 数据库
#### mysql
1. pkg/dbpg (database program)  mysql_* 对mysql client封装
2. 使用 <https://github.com/go-sql-driver/mysql.git> 作为底层实现
#### oracle
### redis
### kafka
1. pkg/kafka 对kafka client封装
2. 使用 <https://github.com/Shopify/sarama.git> 作为底层实现
### zookeeper
