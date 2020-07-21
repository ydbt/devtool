// logger logger.go 文件
package logger

// rs/zerolog 简单封装
// 日志等级描述: panic(5) > fatal(4) > error(3) > warn(2) > info(1) > debug(0) > trace(-1)

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	logger       zerolog.Logger
	funcStrategy FuncStrategy
	iLevel       int8
}

func NewLogger(cfg *LogCfg) *Logger {
	log.Logger = log.Output(
		io.MultiWriter(zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		}, &lumberjack.Logger{
			Filename:   cfg.Path,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackup,
			MaxAge:     cfg.MaxLive,
			Compress:   cfg.IsCompress,
		}),
	)
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix // UNIX Time is faster and smaller than most timestamps
	//zerolog.TimeFieldFormat = zerolog.TimeFieldFormat
	//2006-01-02 15:04:05(自定义格式日期必须使用此时间点) 时间点格式化日期
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.TimestampFieldName = "datetime"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "msg"
	lg := new(Logger)
	lg.SetLevel(cfg.Level)
	lg.SetStrategy(cfg.Strategy)
	lg.logger = log.Logger.With().Str("service", cfg.AppName).Logger()
	lg.Debugf("%s", "global log config init successful")
	return lg
}

// 重新设置日志等级
func (pLog *Logger) SetLevel(level string) {
	lv, err := zerolog.ParseLevel(level)
	if err != nil {
		lv = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(lv)
	pLog.setLevel(int8(lv))
}

// Trace trace 日志等级格式化输出
func (pLog *Logger) Tracef(format string, v ...interface{}) {
	if !pLog.check() {
		return
	}
	callers := pLog.funcStrategy()
	pLog.logger.Trace().Fields(callers).Msgf(format, v...)
}
func (pLog *Logger) Tracemf(fields map[string]interface{}, format string, v ...interface{}) {
	if !pLog.check() {
		return
	}
	callers := pLog.funcStrategy()
	for k, v := range fields {
		callers[k] = v
	}
	pLog.logger.Trace().Fields(callers).Msgf(format, v...)
}

// Debug debug 日志等级格式化输出
func (pLog *Logger) Debugf(format string, v ...interface{}) {
	if !pLog.check() {
		return
	}
	callers := pLog.funcStrategy()
	pLog.logger.Debug().Fields(callers).Msgf(format, v...)
}
func (pLog *Logger) Debugmf(fields map[string]interface{}, format string, v ...interface{}) {
	if !pLog.check() {
		return
	}
	callers := pLog.funcStrategy()
	for k, v := range fields {
		callers[k] = v
	}
	pLog.logger.Debug().Fields(callers).Msgf(format, v...)
}

// LogInfo info 日志等级格式化输出
func (pLog *Logger) Infof(format string, v ...interface{}) {
	if !pLog.check() {
		return
	}
	callers := pLog.funcStrategy()
	pLog.logger.Info().Fields(callers).Msgf(format, v...)
}
func (pLog *Logger) Infomf(fields map[string]interface{}, format string, v ...interface{}) {
	callers := pLog.funcStrategy()
	for k, v := range fields {
		callers[k] = v
	}
	pLog.logger.Info().Fields(callers).Msgf(format, v...)
}

// LogWarn warn 日志等级格式化输出
func (pLog *Logger) Warnf(format string, v ...interface{}) {
	if !pLog.check() {
		return
	}
	callers := pLog.funcStrategy()
	pLog.logger.Warn().Fields(callers).Msgf(format, v...)
}
func (pLog *Logger) Warnmf(fields map[string]interface{}, format string, v ...interface{}) {
	if !pLog.check() {
		return
	}
	callers := pLog.funcStrategy()
	for k, v := range fields {
		callers[k] = v
	}
	pLog.logger.Warn().Fields(callers).Msgf(format, v...)
}

// LogError error 日志等级格式化输出
func (pLog *Logger) Errorf(format string, v ...interface{}) {
	if !pLog.check() {
		return
	}
	callers := pLog.funcStrategy()
	pLog.logger.Error().Fields(callers).Msgf(format, v...)
}
func (pLog *Logger) Errormf(fields map[string]interface{}, format string, v ...interface{}) {
	if !pLog.check() {
		return
	}
	callers := pLog.funcStrategy()
	for k, v := range fields {
		callers[k] = v
	}
	pLog.logger.Error().Fields(callers).Msgf(format, v...)
}

// LogFatal fatal 日志等级格式化输出
func (pLog *Logger) Fatalf(format string, v ...interface{}) {
	if !pLog.check() {
		return
	}
	callers := pLog.funcStrategy()
	pLog.logger.Fatal().Fields(callers).Msgf(format, v...)
}
func (pLog *Logger) Fatalmf(fields map[string]interface{}, format string, v ...interface{}) {
	if !pLog.check() {
		return
	}
	callers := pLog.funcStrategy()
	for k, v := range fields {
		callers[k] = v
	}
	pLog.logger.Fatal().Fields(callers).Msgf(format, v...)
}

func (pLog *Logger) SetStrategy(strategy int) {
	if strategy < 0 {
		strategy = 0
	} else if strategy > 3 {
		strategy = 3
	}
	switch strategy {
	case 0:
		pLog.funcStrategy = strategy0
	case 1:
		pLog.funcStrategy = strategy1
	case 2:
		pLog.funcStrategy = strategy2
	case 3:
		pLog.funcStrategy = strategy3
	default:
		pLog.funcStrategy = strategy0
	}
}

// HotLoadI
// 日志配置热加载
func (pLog *Logger) UpdateCfg(cfg interface{}) {
	if logCfg, ok := cfg.(LogCfg); ok {
		pLog.SetLevel(logCfg.Level)
		pLog.SetStrategy(logCfg.Strategy)
	} else {
		pLog.Warnf("assert LogCfg failed")
	}
}

//
// 日志等级判断
func (pLog *Logger) check() bool {
	return true
}

// setLevel
// 设置日志等级
func (pLog *Logger) setLevel(il int8) {
	pLog.iLevel = il
}
