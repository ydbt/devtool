# kafka client
## 使用
### 示例
#### 生产者
#### 消费者
### josn 字段描述
#### 生产者

|字段|描述|备注|
|-|-|-|
||||
#### 消费者

|字段|类型|描述|参考示例|
|-|-|-|-|
|service_config.brokers|列表|指定kafka机器列表|["nod01:9092","node02:9092"...]|
|service_config.apiversion|字符串|指定kafka服务版本|"1.1.0","0.8.2"|
|custom_config.buffersize|int|消费者客户端缓存区大小||
|groupid|字符串|消费组标识||
|topics|列表|消费组订阅的topic|["topic01","topic02"]|
|goto_offset|int|消费者启动时位置|-1:消费分区最新消息 -2:消费分区可获得的消息(能获得的最旧的消息)|
## 封装
## 测试
### 合法测试
1. 目前仅实现简单合法测试
### 非法测试
### 异常测试
