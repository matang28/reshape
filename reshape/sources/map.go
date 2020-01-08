package sources

import (
	"fmt"
	"github.com/matang28/reshape/reshape"
)

type MapEntry struct {
	Key   interface{}
	Value interface{}
}

type MapSource struct {
	m        map[interface{}]interface{}
	ch       chan interface{}
	isClosed bool
}

func NewMapSource() *MapSource {
	return &MapSource{m: make(map[interface{}]interface{}), ch: make(chan interface{})}
}

func (this *MapSource) Stream() reshape.Stream {
	return reshape.NewStream(this)
}

func (this *MapSource) GetChannel() <-chan interface{} {
	return this.ch
}

func (this *MapSource) Close() error {
	close(this.ch)
	this.isClosed = true
	return nil
}

func (this *MapSource) Put(key, value interface{}) error {
	if !this.isClosed {
		this.m[key] = value
		go func() {
			this.ch <- MapEntry{Key: key, Value: value}
		}()
		return nil
	}
	return fmt.Errorf("cannot put new entries to closed MapSource")
}

func (this *MapSource) Get(key interface{}) (interface{}, bool) {
	value, found := this.m[key]
	return value, found
}
