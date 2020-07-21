package logger_test

import (
	"os"
	"strings"
	"testing"
	"github.com/ydbt/devtool/logger"
)

//
func TestGoroutine(t *testing.T) {
	var testOk = false
	strLogFile := "tmp_test_logoroutine.log"
	strLogContentInKey := "goroutine"
	strLogContentInValue := "pack and pack again"
	strLogContentIn := "i'm here , here are you?"
	strLogContentOut := "where am i?"
	defer func() { // 删除日志文件
		if testOk {
			os.Remove(strLogFile)
		}
	}()
	log := logger.NewLogger(&logger.LogCfg{
		AppName:    "test-app",
		Path:       strLogFile,
		Level:      "debug",
		MaxSize:    100,
		MaxBackup:  100,
		MaxLive:    200,
		IsCompress: false,
		Strategy:   0,
	})
	mfs := make(map[string]interface{})
	mfs[strLogContentInKey] = strLogContentInValue
	lg := logger.NewLogRoutine(log, mfs)
	lg.Tracef("%s", strLogContentIn)
	lg.SetStrategy(-1)
	lg.Debugf("%s", strLogContentIn)
	lg.SetStrategy(0)
	lg.Infof("%s", strLogContentIn)
	lg.SetStrategy(1)
	lg.Warnf("%s", strLogContentIn)
	lg.SetStrategy(2)
	lg.Errorf("%s", strLogContentIn)
	lg.SetStrategy(3)
	lg.Debugf("%s", strLogContentOut)
	defer func() {
		fn, openErr := os.Open(strLogFile)
		if openErr != nil {
			t.Error("open ", strLogFile, " failed!")
			t.Error(openErr)
			return
		}
		defer fn.Close()
		byteLogContent := make([]byte, 2048) // 设置缓冲区大小
		sizeLogContent, readErr := fn.Read(byteLogContent)
		if readErr != nil || sizeLogContent < 1 {
			t.Error(readErr)
			t.Error("log file size=", sizeLogContent)
			return
		}
		strLogContent := string(byteLogContent)
		if strings.Index(strLogContent, strLogContentInKey) == -1 {
			t.Error("write log can't find")
			return
		}
		testOk = true
	}()
}
