package strategies

import (
	"fmt"
	"github.com/matang28/reshape/reshape"
	"time"
)

func tick() {
	time.Sleep(100 * time.Millisecond)
}

func predefinedSource(ch chan interface{}, elements ...interface{}) {
	defer func() {
		recover()
	}()

	for _, e := range elements {
		ch <- e
	}
}

func delayedSource(ch chan interface{}, delay time.Duration, elements ...interface{}) {
	defer func() {
		recover()
	}()

	go func() {
		for _, e := range elements {
			ch <- e
			time.Sleep(delay)
		}
	}()
}

var badTrans reshape.Transformation = func(in interface{}) (out interface{}, err error) {
	return nil, fmt.Errorf("")
}

var plusOneTrans = func(in interface{}) (out interface{}, err error) {
	out = in.(int) + 1
	return out, nil
}

var dropEvens = func(in interface{}) bool {
	return in.(int)%2 != 0
}

type badSink struct {
}

func (this *badSink) Dump(object ...interface{}) error {
	return fmt.Errorf("")
}

func (this *badSink) Close() error {
	return fmt.Errorf("")
}
