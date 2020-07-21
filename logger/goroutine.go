// logger logger.go 文件
package logger

// rs/zerolog 简单封装
// 日志等级描述: panic(5) > fatal(4) > error(3) > warn(2) > info(1) > debug(0) > trace(-1)

type LogRoutine struct {
	log      LogFieldsI             // 日志输出
	exFields map[string]interface{} // 附加信息
}

func NewLogRoutine(logi LogFieldsI, me map[string]interface{}) *LogRoutine {
	lr := new(LogRoutine)
	lr.exFields = make(map[string]interface{})
	for k, v := range me {
		lr.exFields[k] = v
	}
	lr.log = logi
	return lr
}

func (lr *LogRoutine) SetLevel(level string) {
	lr.log.SetLevel(level)
}

func (lr *LogRoutine) SetStrategy(strategy int) {
	lr.log.SetStrategy(strategy)
}

// LogTrace trace 日志等级格式化输出
func (lr *LogRoutine) Tracef(format string, v ...interface{}) {
	lr.log.Tracemf(lr.exFields, format, v...)
}

// LogDebug debug 日志等级格式化输出
func (lr *LogRoutine) Debugf(format string, v ...interface{}) {
	lr.log.Debugmf(lr.exFields, format, v...)
}

// LogInfo info 日志等级格式化输出
func (lr *LogRoutine) Infof(format string, v ...interface{}) {
	lr.log.Infomf(lr.exFields, format, v...)
}

// LogWarn warn 日志等级格式化输出
func (lr *LogRoutine) Warnf(format string, v ...interface{}) {
	lr.log.Warnmf(lr.exFields, format, v...)
}

// LogError error 日志等级格式化输出
func (lr *LogRoutine) Errorf(format string, v ...interface{}) {
	lr.log.Errormf(lr.exFields, format, v...)
}

// LogFatal fatal 日志等级格式化输出
func (lr *LogRoutine) Fatalf(format string, v ...interface{}) {
	lr.log.Fatalmf(lr.exFields, format, v...)
}
