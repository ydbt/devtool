package logger_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"github.com/ydbt/devtool/v3/logger"
)

// 回滚日志文件测试
func TestLogRollback(_t *testing.T) {
	var testOk = false
	strLogFile := "tmp_test_rollback.log"
	strLogContentIn := "hi! i'm here, where are you?"
	strLogContentOut := "where am i?"
	defer func() { // 删除日志文件
		if testOk {
			os.Remove(strLogFile)
		}
	}()
	{ // 日志写入操作
		cfgJS := fmt.Sprintf("{\"appname\":\"test-app\",\"path\":\"%s\",\"level\":\"trace\",\"maxsize\":100,\"maxbackup\":100,\"maxlive\":60,\"compress\":false,\"strategy\":2}", strLogFile)
		cfgLog, _ := logger.Json2LogCfg(cfgJS)
		logRollback := logger.NewLogger(cfgLog)
		logRollback.SetStrategy(0)
		logRollback.Tracef(strLogContentIn)
		logRollback.SetStrategy(1)
		logRollback.Debugf(strLogContentIn)
		logRollback.SetLevel("xxxx") // trace 日志等级
		logRollback.Infof(strLogContentIn)
		logRollback.SetStrategy(2)
		logRollback.Warnf(strLogContentIn)
		logRollback.SetStrategy(3)
		logRollback.Errorf(strLogContentIn)
		//		logRollback.Fatalf(strLogContentIn)
		logRollback.SetLevel("info") // trace 日志等级
		logRollback.Debugf(strLogContentOut)
	}
	defer func() {
		fn, openErr := os.Open(strLogFile)
		if openErr != nil {
			_t.Error("open ", strLogFile, " failed!")
			_t.Error(openErr)
			return
		}
		defer fn.Close()
		byteLogContent := make([]byte, 1024) // 设置缓冲区大小
		sizeLogContent, readErr := fn.Read(byteLogContent)
		if readErr != nil || sizeLogContent < 1 {
			_t.Error(readErr)
			_t.Error("log file size=", sizeLogContent)
			return
		}
		strLogContent := string(byteLogContent)
		if strings.Index(strLogContent, strLogContentIn) == -1 {
			_t.Error("write log can't find")
			return
		}
		if strings.Index(strLogContent, strLogContentOut) != -1 {
			_t.Error("set level failed, low level log exist")
			return
		}
		testOk = true
	}()
}
