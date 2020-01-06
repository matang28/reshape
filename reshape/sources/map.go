package sources

import (
	"github.com/matang28/reshape/reshape"
)

type MapEntry struct {
	Key   interface{}
	Value interface{}
}

type MapSource struct {
	m  map[interface{}]interface{}
	ch chan interface{}
}

func NewMapSource() *MapSource {
	return &MapSource{m: make(map[interface{}]interface{}), ch: make(chan interface{})}
}

func (this *MapSource) Stream() reshape.Stream {
	return reshape.NewStream(this.ch)
}

func (this *MapSource) CloseGracefully() error {
	close(this.ch)
	return nil
}

func (this *MapSource) Put(key, value interface{}) {
	this.m[key] = value
	this.ch <- MapEntry{Key: key, Value: value}
}

func (this *MapSource) Get(key interface{}) (interface{}, bool) {
	value, found := this.m[key]
	return value, found
}
