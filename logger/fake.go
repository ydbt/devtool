// 测试日志使用接口
package logger

import "fmt"

// rs/zerolog 简单封装
// 日志等级描述: panic(5) > fatal(4) > error(3) > warn(2) > info(1) > debug(0) > trace(-1)

type LoggerFake struct {
	IsConsole bool
	Level     string
}

func (pLog *LoggerFake) SetLevel(level string) {
	pLog.Level = level
}
func (pLog *LoggerFake) Init(_bConsole bool) {
	pLog.IsConsole = _bConsole
}

func (pLog *LoggerFake) LogConsole(format string, v ...interface{}) {
	if pLog.IsConsole {
		fmt.Print(pLog.Level, ", ")
		if v == nil || len(v) == 0 {
			fmt.Println(format)
		} else {
			fmt.Printf(format, v...)
			fmt.Println("")
		}
		fmt.Println("")
	}
}
func (pLog *LoggerFake) SetStrategy(strategy int) {

}

// LogTrace trace 日志等级格式化输出
func (pLog *LoggerFake) Tracef(format string, v ...interface{}) {
	pLog.LogConsole(format, v...)
}

// LogDebug debug 日志等级格式化输出
func (pLog *LoggerFake) Debugf(format string, v ...interface{}) {
	pLog.LogConsole(format, v...)
}

// LogInfo info 日志等级格式化输出
func (pLog *LoggerFake) Infof(format string, v ...interface{}) {
	pLog.LogConsole(format, v...)
}

// LogWarn warn 日志等级格式化输出
func (pLog *LoggerFake) Warnf(format string, v ...interface{}) {
	pLog.LogConsole(format, v...)
}

// LogError error 日志等级格式化输出
func (pLog *LoggerFake) Errorf(format string, v ...interface{}) {
	pLog.LogConsole(format, v...)
}

// LogFatal fatal 日志等级格式化输出
func (pLog *LoggerFake) Fatalf(format string, v ...interface{}) {
	pLog.LogConsole(format, v...)
}

// LogFatal fatal 日志等级格式化输出
func (pLog *LoggerFake) Panicf(format string, v ...interface{}) {
	pLog.LogConsole(format, v...)
}
