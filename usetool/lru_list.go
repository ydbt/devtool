package usetool

import "sync"

type Element struct {
	next, prev *Element    // 执行相邻元素的指针
	key        string      // 快速访问key
	Value      interface{} // 存储的值
}

// LruList
// 循环列表
// 参考 "container/list"
type LruList struct {
	root      Element
	length    int
	fastTable map[string]*Element
	mtxElem   sync.Mutex
}

func NewLruList() *LruList {
	ll := new(LruList)
	ll.Clear()
	return ll
}

func (ll *LruList) Clear() {
	ll.mtxElem.Lock()
	defer ll.mtxElem.Unlock()
	ll.length = 0
	ll.fastTable = make(map[string]*Element)
	ll.root.prev = &ll.root
	ll.root.next = &ll.root
}

func (ll *LruList) Len() int {
	ll.mtxElem.Lock()
	defer ll.mtxElem.Unlock()
	return ll.length
}

func (ll *LruList) Front() *Element {
	ll.mtxElem.Lock()
	defer ll.mtxElem.Unlock()
	if ll.length == 0 {
		return nil
	}
	return ll.root.next
}

func (ll *LruList) Back() *Element {
	ll.mtxElem.Lock()
	defer ll.mtxElem.Unlock()
	if ll.length == 0 {
		return nil
	}
	return ll.root.prev
}

// Index
// 快速访问
func (ll *LruList) Index(key string) *Element {
	ll.mtxElem.Lock()
	defer ll.mtxElem.Unlock()
	if e, ok := ll.fastTable[key]; ok {
		return e
	} else {
		return nil
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (ll *LruList) insert(e, at *Element) *Element {
	ll.mtxElem.Lock()
	defer ll.mtxElem.Unlock()
	if v, ok := ll.fastTable[e.key]; ok {
		v.Value = e.Value
	} else {
		n := at.next
		at.next = e
		e.prev = at
		e.next = n
		n.prev = e
		ll.length++
		ll.fastTable[e.key] = e
	}
	return e
}

// Pushback
// 尾部插入数据
func (ll *LruList) Pushback(key string, v interface{}) *Element {
	e := &Element{Value: v,
		key: key,
	}
	return ll.insert(e, ll.root.prev)
}

// Pushfront
// 头部插入数据
func (ll *LruList) Pushfront(key string, v interface{}) *Element {
	e := &Element{Value: v,
		key: key,
	}
	return ll.insert(e, &ll.root)
}

// Remove
// 删除节点
func (ll *LruList) Remove(key string) *Element {
	ll.mtxElem.Lock()
	defer ll.mtxElem.Unlock()
	e, ok := ll.fastTable[key]
	if !ok {
		return nil
	}
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil
	e.prev = nil
	ll.length--
	delete(ll.fastTable, key)
	return e
}

// Popback
// 移除尾部元素
func (ll *LruList) Popback() *Element {
	if ll.length > 0 {
		e := ll.root.prev
		e.next = &ll.root
		ll.root.prev = e.prev
		e.next = nil
		e.prev = nil
		return e
	}
	return nil
}

// Popfront
// 移除首部元素
func (ll *LruList) Popfront() *Element {
	if ll.length > 0 {
		e := ll.root.next
		e.next.prev = &ll.root
		ll.root.next = e.next
		e.next = nil
		e.prev = nil
		return e
	}
	return nil
}

// Move2front
// 将元素移动到首部
func (ll *LruList) Move2front(key string) {
	ll.mtxElem.Lock()
	defer ll.mtxElem.Unlock()
	e, ok := ll.fastTable[key]
	if !ok {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.next = ll.root.next
	e.prev = &ll.root
	ll.root.next = e

}

// Move2back
// 将元素移动到尾部
func (ll *LruList) Move2back(key string) {
	ll.mtxElem.Lock()
	defer ll.mtxElem.Unlock()
	e, ok := ll.fastTable[key]
	if !ok {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = ll.root.prev
	e.next = &ll.root
	ll.root.prev = e
}

// TraverseValue
// 遍历插入的值
func (ll *LruList) TraverseValue() []interface{} {
	ll.mtxElem.Lock()
	defer ll.mtxElem.Unlock()
	lv := make([]interface{}, ll.length)
	beg := ll.root.next
	end := &ll.root
	i := 0
	for beg != end {
		lv[i] = beg.Value
		beg = beg.next
		i++
	}
	return lv
}

// TraverseElement
// 遍历节点
func (ll *LruList) TraverseElement() []*Element {
	ll.mtxElem.Lock()
	defer ll.mtxElem.Unlock()
	le := make([]*Element, ll.length)
	beg := ll.root.next
	end := &ll.root
	i := 0
	for beg != end {
		le[i] = beg
		beg = beg.next
		i++
	}
	return le
}

// BIndex
// Balance Index 当key被访问时，将节点移动到首部
func (ll *LruList) BIndex(key string) *Element {
	ll.mtxElem.Lock()
	defer ll.mtxElem.Unlock()
	v, ok := ll.fastTable[key]
	if !ok {
		return nil
	}
	v.prev.next = v.next
	v.next.prev = v.prev

	v.next = ll.root.next
	ll.root.next = v
	return v
}
