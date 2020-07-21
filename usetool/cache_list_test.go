package usetool_test

import (
	"fmt"
	"testing"
	usetool "github.com/ydbt/devtool/v3/usetool"
)

func TestCachedListContent(t *testing.T) {
	cl := usetool.NewCacheList(0)
	cl.Push("000", 0)
	//	t.Log(cl)
	cl.Push("001", 1)
	//	t.Log(cl)
	cl.Push("002", 2)
	//	t.Log(cl)
	cl.Push("003", 3)
	//	t.Log(cl)
	iarrActual := cl.TraverseValue()
	iarrExpect := []int{3, 2, 1, 0}
	for i, v := range iarrExpect {
		if iarrActual[i].(int) != v {
			t.Error("push put 'dat' head")
		}
	}
}

func TestCachedListPush(t *testing.T) {
	cl := usetool.NewCacheList(50)
	for i := 0; i < 50; i++ {
		k := fmt.Sprintf("%03d", i)
		cl.Push(k, i)
	}
	k000 := "000"
	k049 := "049"
	e000 := cl.Index(k000)
	e049 := cl.Index(k049)
	if e000 != nil {
		v000 := e000.Value.(int)
		if v000 != 0 {
			t.Errorf("push(000->049) ['000']=%d", v000)
		}
	}
	if e049 != nil {
		v049 := e049.Value.(int)
		if v049 != 49 {
			t.Errorf("push(000->049) ['049']=%d", v049)
		}
	}
	cl.Push("050", 50) // 此时51个元素 51 > 50
	cl.Push("051", 51) // 此时52个元素
	e000 = cl.Index(k000)
	if e000 != nil {
		t.Error("push 000->051 ['000'] should removed , please check")
		t.Log(cl)
	}
}

func TestCachedListBIndex(t *testing.T) {
	cl := usetool.NewCacheList(50)
	for i := 0; i < 50; i++ {
		k := fmt.Sprintf("%03d", i)
		cl.Push(k, i)
	}
	k000 := "000"
	k049 := "049"
	e000 := cl.Index(k000)
	e049 := cl.Index(k049)
	if e000 != nil {
		v000 := e000.Value.(int)
		if v000 != 0 {
			t.Errorf("push(000->049) ['000']=%d", v000)
		}
	} else {
		t.Error("push(000->049) buffer doesn't full")
		t.Log(cl)
	}
	if e049 != nil {
		v049 := e049.Value.(int)
		if v049 != 49 {
			t.Errorf("push(000->049) ['049']=%d", v049)
		}
	} else {
		t.Error("push(000->049) buffer doesn't full")
		t.Log(cl)
	}
	// 上述为构建满缓冲池
	e000 = cl.BIndex(k000)
	k001 := "001"
	e001 := cl.BIndex(k001)
	cl.Push("050", 50)
	cl.Push("051", 51)
	e000 = cl.Index(k000)
	if e000 == nil {
		t.Error("cachedlist balanceindex failed")
	}
	e001 = cl.Index(k001)
	if e001 == nil {
		t.Error("cachedlist balanceindex failed")
	}
	e002 := cl.Index("002")
	if e002 != nil {
		t.Error("remove cached data failed, when buffer be fulled")
	}
	//t.Log(cl.TraverseValue())
}
