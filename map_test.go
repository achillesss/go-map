package gomap

import (
	"testing"

	"github.com/achillesss/log"
)

func TestMapInt(t *testing.T) {
	srcMap := make(map[int]int)
	// var srcMap map[int]int
	m := NewMap(srcMap)
	go m.Handler()
	m.Add(1, 1)
	var q int
	m.Query(1, &q)
	if q != 1 {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
	m.Delete(1)
	m.Query(1, &q)
	if q != 0 {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
	m.Close()
}

func TestMapString(t *testing.T) {
	var srcMap map[string]string
	m := NewMap(srcMap)
	go m.Handler()
	m.Add("1", "1")
	var q string
	m.Query("1", &q)
	if q != "1" {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
	m.Delete("1")
	m.Query("1", &q)
	if q != "" {
		t.Errorf("%s failed.", log.FuncName())
		return
	}

}

func TestMapValueStruct(t *testing.T) {
	type A struct {
		a bool
		b int
		c string
		d []int
	}
	var srcMap map[int]*A
	var a A
	a.a = true
	a.b = 2
	a.c = "hello"
	a.d = []int{3, 4}

	m := NewMap(srcMap)
	go m.Handler()
	m.Add(9, &a)
	var q *A
	m.Query(9, &q)
	if q == nil || q.a != a.a || q.b != a.b && q.c != a.c || (len(q.d) != len(q.d)) || q.d[0] != q.d[0] || q.d[1] != q.d[1] {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
	m.Delete(0)
	m.Query(0, &q)
	if q != nil {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
}
