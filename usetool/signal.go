package usetool

import (
	"os"
	"sync"
)

func NewProcessSignaler() *ProcessSignaler {
	psr := new(ProcessSignaler)
	psr.subscribers = make(map[os.Signal]*[]RoutineSignal)
	return psr
}

func (psr *ProcessSignaler) WaitSignalProcess(sg os.Signal) {
	psr.mtx.Lock()
	defer psr.mtx.Unlock()
	if subers, ok := psr.subscribers[sg]; ok {
		for _, v := range *subers {
			v <- struct{}{}
		}
	}
}

// TopicOSSignal
// 订阅指定的进程信号
func (psr *ProcessSignaler) TopicOSSignal(sg os.Signal) RoutineSignal {
	psr.mtx.Lock()
	defer psr.mtx.Unlock()
	tsg := make(RoutineSignal, 1)
	if lsg, ok := psr.subscribers[sg]; ok {
		*lsg = append(*lsg, tsg)
	} else {
		lsg := make([]RoutineSignal, 1)
		psr.subscribers[sg] = &lsg
		lsg[0] = tsg
	}
	return tsg
}

// RoutineSignal
// 进行关闭协程退出信号
type RoutineSignal chan struct{}

// ProcessSignaler
// 进程信号管理者
type ProcessSignaler struct {
	subscribers map[os.Signal]*[]RoutineSignal
	mtx         sync.Mutex
}

//
