package sinks

import "fmt"

type KeyExtractor = func(in interface{}) interface{}

type MapSink struct {
	keyExtractor KeyExtractor
	m            map[interface{}]interface{}
}

func NewMapSink(keyExtractor KeyExtractor) *MapSink {
	return &MapSink{keyExtractor: keyExtractor, m: make(map[interface{}]interface{})}
}

func (this *MapSink) Dump(object ...interface{}) error {
	for _, elem := range object {
		key := this.keyExtractor(elem)
		if key == nil {
			return fmt.Errorf("cannot extract key from %+v", elem)
		}
		this.m[key] = elem
	}
	return nil
}

func (this *MapSink) Close() error {
	return nil
}

func (this *MapSink) Get() map[interface{}]interface{} {
	return this.m
}
