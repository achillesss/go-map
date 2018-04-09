package gomap

import (
	"fmt"
	"reflect"
)

type GoMap struct {
	instance          interface{}
	addChan           chan map[interface{}]interface{}
	delChan           chan interface{}
	queryChan         chan interface{}
	queryRespChan     chan map[interface{}]interface{}
	interfaceChan     chan struct{}
	interfaceRespChan chan interface{}
	dropChan          chan struct{}
}

// NewMap creates a map
func NewMap(srcMap interface{}) *GoMap {
	srcType := reflect.TypeOf(srcMap)
	if srcType.Kind() != reflect.Map {
		panic("src not map")
	}
	var m GoMap
	m.instance = srcMap
	m.addChan = make(chan map[interface{}]interface{})
	m.delChan = make(chan interface{})
	m.queryChan = make(chan interface{})
	m.queryRespChan = make(chan map[interface{}]interface{})
	m.interfaceChan = make(chan struct{})
	m.interfaceRespChan = make(chan interface{})
	m.dropChan = make(chan struct{})
	return &m
}

func isMap(src interface{}) bool {
	return reflect.TypeOf(src).Kind() == reflect.Map
}

// MapHandler handles map
func (gm GoMap) Handler() {
	mapValue := reflect.ValueOf(gm.instance)
	mapType := reflect.TypeOf(gm.instance)

	if mapValue.IsNil() {
		mapValue = reflect.MakeMap(mapType)
	}
	keysType := mapType.Key()
	valueType := mapType.Elem()

	for {
		select {

		// add
		case m := <-gm.addChan:
			for k, v := range m {
				kv := reflect.ValueOf(k)
				vv := reflect.ValueOf(v)
				if vv.IsValid() {
					mapValue.SetMapIndex(kv, vv)
				}
			}

		// delete
		case m := <-gm.delChan:
			mv := reflect.ValueOf(m)
			mapValue.SetMapIndex(mv, reflect.Value{})

		// query
		case m := <-gm.queryChan:
			kt := reflect.TypeOf(m)
			if kt.Kind() != keysType.Kind() {
				gm.queryRespChan <- nil
				continue
			}

			kv := reflect.ValueOf(m)
			v := mapValue.MapIndex(kv)
			newV := reflect.New(valueType)

			if !v.IsValid() {
				newV.Elem().Set(reflect.Zero(valueType))
			} else {
				newV.Elem().Set(v)
			}
			gm.queryRespChan <- map[interface{}]interface{}{m: newV.Elem().Interface()}

		// change to interface{}
		case <-gm.interfaceChan:
			gm.interfaceRespChan <- mapValue.Interface()
		// drop, interrupt select loop
		case <-gm.dropChan:
			return
		}

	}
}

// Add add key: value to map
func (gm *GoMap) Add(key, value interface{}) {
	gm.addChan <- map[interface{}]interface{}{key: value}
}

func (gm *GoMap) Delete(key interface{}) {
	gm.delChan <- key
}

func (gm *GoMap) pickQueryResp(key interface{}) interface{} {
	for resp := range gm.queryRespChan {
		if resp == nil {
			break
		}
		for k, v := range resp {
			if reflect.DeepEqual(key, k) {
				return v
			}
		}
		gm.queryRespChan <- resp
	}
	return nil
}

func (gm *GoMap) Query(key interface{}, dst interface{}) error {
	gm.queryChan <- key
	v := gm.pickQueryResp(key)
	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		return fmt.Errorf("bad dst type")
	}

	dv := reflect.ValueOf(dst)
	if dv.Kind() != reflect.Ptr {
		panic("dst not pointer")
	}

	dv.Elem().Set(reflect.ValueOf(v))
	return nil
}

func (gm *GoMap) Interface() interface{} {
	gm.interfaceChan <- struct{}{}
	return <-gm.interfaceRespChan
}

// Close means no other coming actions after it.
func (gm *GoMap) Close() {
	gm.dropChan <- struct{}{}
}
