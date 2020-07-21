package usetool

import "os"

// ProcessSignalI
// 进行信号订阅
type ProcessSignalI interface {
	TopicOSSignal(os.Signal) RoutineSignal
}
