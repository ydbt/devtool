package usetool

import "sync"

// CacheList
// 数据缓存
type CacheList struct {
	root      Element
	length    int
	maxSize   int
	fastTable map[string]*Element
	mtxElem   sync.Mutex
}

func NewCacheList(maxCap int) *CacheList {
	cl := new(CacheList)
	if maxCap < 30 {
		cl.maxSize = 30
	} else if maxCap > 5000 {
		cl.maxSize = 5000
	} else {
		cl.maxSize = maxCap
	}

	cl.Clear()
	return cl
}

func (cl *CacheList) Clear() {
	cl.mtxElem.Lock()
	defer cl.mtxElem.Unlock()
	cl.length = 0
	cl.fastTable = make(map[string]*Element)
	cl.root.prev = &cl.root
	cl.root.next = &cl.root
}

// Resize
// 调整缓存大小
func (cl *CacheList) Resize(maxCap int) {
	cl.mtxElem.Lock()
	defer cl.mtxElem.Unlock()
	if maxCap < 30 {
		cl.maxSize = 30
	} else if maxCap > 5000 {
		cl.maxSize = 5000
	} else {
		cl.maxSize = maxCap
	}
	if cl.length < cl.maxSize {
		return
	}

	// 当缓存容量不够时删除最后的 1/5 数据 -- 策略可调
	delSize := cl.maxSize/5 + (cl.length - cl.maxSize)
	delLast := cl.root.prev
	for i := 0; i < delSize; i++ {
		de := delLast
		delLast = delLast.prev
		delete(cl.fastTable, de.key)
	}
	cl.length = cl.length - delSize
	cl.root.prev = delLast
	delLast.next = &cl.root
}

// Index
// 快速访问
func (cl *CacheList) Index(key string) *Element {
	cl.mtxElem.Lock()
	defer cl.mtxElem.Unlock()
	if e, ok := cl.fastTable[key]; ok {
		return e
	} else {
		return nil
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (cl *CacheList) insert(e, at *Element) *Element {
	cl.mtxElem.Lock()
	defer cl.mtxElem.Unlock()
	if cl.length > cl.maxSize {
		// 当缓存容量不够时删除最后的 1/5 数据 -- 策略可调
		delSize := cl.maxSize / 5
		delLast := cl.root.prev
		for i := 0; i < delSize; i++ {
			de := delLast
			delLast = delLast.prev
			delete(cl.fastTable, de.key)
		}
		cl.length = cl.length - delSize
		cl.root.prev = delLast
		delLast.next = &cl.root
	}
	if v, ok := cl.fastTable[e.key]; ok {
		v.Value = e.Value
	} else {
		n := at.next
		at.next = e
		e.prev = at
		e.next = n
		n.prev = e
		cl.length++
		cl.fastTable[e.key] = e
	}
	return e
}

// Pushfront
// 头部插入数据
func (cl *CacheList) Push(key string, v interface{}) *Element {
	e := &Element{Value: v,
		key: key,
	}
	return cl.insert(e, &cl.root)
}

// TraverseValue
// 遍历插入的值
func (cl *CacheList) TraverseValue() []interface{} {
	cl.mtxElem.Lock()
	defer cl.mtxElem.Unlock()
	lv := make([]interface{}, cl.length)
	beg := cl.root.next
	end := &cl.root
	i := 0
	for beg != end {
		lv[i] = beg.Value
		beg = beg.next
		i++
	}
	return lv
}

// BIndex
// Balance Index 当key被访问时，将节点移动到首部
func (cl *CacheList) BIndex(key string) *Element {
	cl.mtxElem.Lock()
	defer cl.mtxElem.Unlock()
	v, ok := cl.fastTable[key]
	if !ok {
		return nil
	}
	v.prev.next = v.next
	v.next.prev = v.prev

	v.next = cl.root.next
	cl.root.next = v
	return v
}
