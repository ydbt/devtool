package logger

type LogFieldsI interface {
	SetLevel(string)
	SetStrategy(int)
	Tracemf(fields map[string]interface{}, forma string, args ...interface{})
	Debugmf(fields map[string]interface{}, format string, args ...interface{})
	Infomf(fields map[string]interface{}, format string, args ...interface{})
	Warnmf(fields map[string]interface{}, format string, args ...interface{})
	Errormf(fields map[string]interface{}, format string, args ...interface{})
	Fatalmf(fields map[string]interface{}, format string, args ...interface{})
}

type LogI interface {
	SetLevel(string)
	SetStrategy(int)
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}
