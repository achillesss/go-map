package gomap

import (
	"fmt"
	"reflect"
)

type goMap reflect.Value

// NewMap creates a map
func NewMap(key, value interface{}) goMap {
	mt := reflect.MapOf(reflect.TypeOf(key), reflect.TypeOf(value))
	return goMap(reflect.New(mt))
}

func isMap(src interface{}) bool {
	return reflect.TypeOf(src).Kind() == reflect.Map
}

// MapHandler handles map
func (gm goMap) Handler() {
	mapValue := reflect.ValueOf(gm)
	mapType := reflect.TypeOf(gm)
	if mapValue.IsNil() {
		mapValue = reflect.New(mapType)
	}
	mapKeys := mapValue.MapKeys()
	keysType := reflect.TypeOf(mapKeys).Elem()

	for {
		select {

		// add
		case m := <-mapAddChan:
			for k, v := range m {
				kv := reflect.ValueOf(k)
				vv := reflect.ValueOf(v)
				if vv.IsValid() {
					mapValue.SetMapIndex(kv, vv)
				}
			}

		// delete
		case m := <-mapDelChan:
			zeroValue := reflect.New(mapValue.Elem().Type())
			mv := reflect.ValueOf(m)
			mapValue.SetMapIndex(mv, zeroValue)

		// query
		case m := <-mapQueryChan:
			kt := reflect.TypeOf(m)
			if kt != keysType {
				mapQueryRespChan <- nil
				continue
			}
			kv := reflect.ValueOf(m)
			v := mapValue.MapIndex(kv)
			mapQueryRespChan <- map[interface{}]interface{}{m: v}
		}

	}
}

var mapAddChan chan map[interface{}]interface{}

// Add add key: value to map
func Add(key, value interface{}) {
	mapAddChan <- map[interface{}]interface{}{key: value}
}

var mapDelChan chan interface{}

func Delete(key interface{}) {
	mapDelChan <- key
}

var mapQueryChan chan interface{}
var mapQueryRespChan chan map[interface{}]interface{}

func pickQueryResp(key interface{}) interface{} {
	for resp := range mapQueryRespChan {
		if resp == nil {
			break
		}
		for k, v := range resp {
			if reflect.DeepEqual(key, k) {
				return v
			}
		}
		mapQueryRespChan <- resp
	}
	return nil
}

func Query(key interface{}, dst interface{}) error {
	mapQueryChan <- key
	v := pickQueryResp(key)
	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		return fmt.Errorf("bad dst type")
	}
	dv := reflect.ValueOf(dst)
	dv.Set(reflect.ValueOf(v))
	return nil
}

func init() {
	mapAddChan = make(chan map[interface{}]interface{})
	mapDelChan = make(chan interface{})
	mapQueryChan = make(chan interface{})
	mapQueryRespChan = make(chan map[interface{}]interface{})
}
