package usetool_test

import (
	"testing"
	usetool "gitlab.qn.com/ydbt/usetool"
)

// 缓存快表

// TestListNew
// 空快表查询
func TestListNew(t *testing.T) {
	list := usetool.NewLruList()
	if l := list.Len(); l != 0 {
		t.Errorf("just new NewLruList object , Len() != 0(%d)", l)
	}
	if v := list.Back(); v != nil {
		t.Errorf("just new NewLruList object , Back() != nil(%v)", v)
	}
	if v := list.Front(); v != nil {
		t.Errorf("just new NewLruList object , Front() != nil(%v)", v)
	}
}

// 插入数据测试
func TestListPush(t *testing.T) {
	list := usetool.NewLruList()
	list.Pushback("1", 1)
	list.Pushfront("2", 2)
	if l := list.Len(); l != 2 {
		t.Errorf("list push two data , Len() != 2(%d)", l)
	}
	v := list.Index("1")
	v1 := v.Value.(int)
	if v1 != 1 {
		t.Errorf("pushback(\"1\",1) & pushfront(\"2\",2) , inex(\"1\") != 1(%d)", v1)
	}
	list.Pushfront("3", 3)
	v = list.Front()
	v1 = v.Value.(int)
	if v1 != 3 {
		t.Errorf("pushfront(\"3\",3), but front().Value != 3 (%d)", v1)
	}
}

// 删除节点测试
func TestListRemove(t *testing.T) {
	list := usetool.NewLruList()
	list.Pushback("1", 1)
	list.Pushback("2", 2)
	list.Pushback("3", 3)
	list.Pushback("4", 4)
	list.Pushback("5", 5)
	if l := list.Len(); l != 5 {
		t.Errorf("pushback(\"1\",\"2\",\"3\",\"4\",\"5\") , Len != 5(%d)", l)
	}
	le := list.TraverseElement()
	for _, tv := range le {
		t.Log("before:", tv)
	}
	t.Log(list)
	v := list.Remove("3")
	if v.Value.(int) != 3 {
		t.Error("remove exception")
	}
	if l := list.Len(); l != 4 {
		t.Error("remove data , length doesn't descrise")
	}
	v = list.Index("3")
	if v != nil {
		t.Error("data be removed, should index")
	}
	lv := list.TraverseValue()
	for _, tv := range lv {
		t.Log("after:", tv)
	}

}
