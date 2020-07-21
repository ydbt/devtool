package redisclient_test

import (
	"testing"
	"time"

	"github.com/ydbt/devtool/v1/logger"
	"github.com/ydbt/devtool/v1/redisclient"
)

var g_tc *redisclient.RedisClient
func init(){
	g_tc = createTestRedisClient()
}

func createTestRedisClient() *redisclient.RedisClient {
	cfg := &redisclient.RedisCfg{
		Addrs:        []string{"10.130.29.88:7000","10.130.29.88:7001","10.130.29.88:7002","10.130.29.88:7003","10.130.29.88:7004","10.130.29.88:7005"},
		User:         "",
		Password:     "",
		Database:     0,
		MaxRetry:     -1,
		PoolSize:     10,
		MinIdleConns: 5,
		PoolTimeout:  60000,
		IdleTimeout:  60000,
	}
	lg := &logger.LoggerFake{
		IsConsole: true,
		Level: "trace",
	}
	rc := redisclient.NewRedisClient(cfg,lg)
	return rc
}


/* --------------------------- redis 字符串接口测试 ------------------------------ */
func TestStringSet(t *testing.T) {
	key := "ut_key_set_0"
	expect := "2020-07-02"
	err := g_tc.Set(key , expect,time.Second * 10)
	if err != nil {
		t.Error(err)
		return
	}
	actaul,err := g_tc.Get(key)
	if err != nil {
		t.Error(err)
		return
	}
	if actaul != expect {
		t.Errorf("test logic error: \"%s\" != \"%s\"",expect,actaul)
	}
}

func TestKeyScan(t *testing.T) {
	key1 := "ut_key_scan_1"
	key2 := "ut_key_scan_2"
	expect := "refactor redis client"
	actual := ""
	err := g_tc.Set(key1,expect,time.Second * 60)
	if err != nil {
		t.Error(err)
		return
	}
	err = g_tc.Set(key2,expect,time.Second * 60)
	if err != nil {
		t.Error(err)
		return
	}
	actual,err = g_tc.Get(key1)
	if actual != expect {
		t.Errorf("%s != %s",expect,actual)
		return
	}
	var exprKeys []string
	var keys []string
	nextCursor := uint64(0)
	for true {
		keys,nextCursor,err = g_tc.Scan(nextCursor,"ut*",100)
		exprKeys = append(exprKeys,keys...)
		if err != nil {
			t.Error(err)
			return
		}
		if nextCursor == 0 {
			break
		}
	}
}

func TestKeyExpire(t *testing.T) {
	key := "ut_key_expire_0"
	val := "refactor redis client"
	err := g_tc.Set(key,val,time.Second * 10)
	if err != nil {
		t.Error(err)
		return
	}
	status,err := g_tc.Pexpire(key,time.Second * 2)
	if err != nil {
		t.Error(err)
		return
	}
	if ! status {
		t.Errorf("maybe %s not exist in redis",key)
		return
	}
	alive,err := g_tc.Pttl(key)
	if err != nil {
		t.Error(err)
		return
	}
	if alive < time.Second {
		t.Error("maybe error , i don't think raise by latency")
	}
}







