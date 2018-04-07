package gomap

import "testing"
import "github.com/achillesss/log"

func TestMapInt(t *testing.T) {
	// srcMap := make(map[int]int)
	var srcMap map[int]int
	m := NewMap(srcMap)
	go m.Handler()
	m.Add(1, 1)
	m.Add(2, 2)
	m.Add(3, 3)
	var q int
	m.Query(1, &q)
	if q != 1 {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
	m.Query(2, &q)
	if q != 2 {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
	m.Query(3, &q)
	if q != 3 {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
}

func TestMapString(t *testing.T) {
	// srcMap := make(map[int]int)
	var srcMap map[string]string
	m := NewMap(srcMap)
	go m.Handler()
	m.Add("1", "1")
	m.Add("2", "2")
	m.Add("3", "3")
	var q string
	m.Query("1", &q)
	if q != "1" {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
	m.Query("2", &q)
	if q != "2" {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
	m.Query("3", &q)
	if q != "3" {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
}
