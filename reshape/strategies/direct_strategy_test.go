package strategies

import (
	"github.com/matang28/reshape/reshape"
	"github.com/matang28/reshape/reshape/sinks"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestDirectStrategy_Solve_HappyCase(t *testing.T) {
	strg := directStrategy{}
	src := make(chan interface{})
	sink := sinks.NewArraySink()

	go strg.Solve(src, []interface{}{plusOneTrans, plusOneTrans, dropEvens, sink})

	predefinedSource(src, -1)
	tick()
	assert.EqualValues(t, []interface{}{1}, sink.Get())

	predefinedSource(src, 1, 2, 3, 4, 5)
	tick()
	assert.EqualValues(t, []interface{}{1, 3, 5, 7}, sink.Get())

	predefinedSource(src, 5, 11, 10)
	tick()
	assert.EqualValues(t, []interface{}{1, 3, 5, 7, 7, 13}, sink.Get())

	close(src)
}

func TestDirectStrategy_Solve_BadTransformation(t *testing.T) {
	strg := directStrategy{}
	src := make(chan interface{})
	sink := sinks.NewArraySink()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := strg.Solve(src, []interface{}{plusOneTrans, plusOneTrans, badTrans, sink})
		close(src)
		assert.NotNil(t, err)
		wg.Done()
	}()

	predefinedSource(src, 1, 2, 3, 4)
	tick()
	assert.Nil(t, sink.Get())
	wg.Wait()
}

func TestDirectStrategy_Solve_BadSink(t *testing.T) {
	strg := directStrategy{}
	src := make(chan interface{})
	sink := &badSink{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := strg.Solve(src, []interface{}{plusOneTrans, plusOneTrans, sink})
		close(src)
		assert.NotNil(t, err)
		wg.Done()
	}()

	predefinedSource(src, 1, 2, 3, 4)
	tick()
	wg.Wait()
}

func TestDirectStrategy_Solve_UnrecognizedHandler(t *testing.T) {
	strg := directStrategy{}
	src := make(chan interface{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := strg.Solve(src, []interface{}{plusOneTrans, plusOneTrans, 10})
		close(src)
		assert.NotNil(t, err)
		_, ok := err.(*reshape.UnrecognizedHandler)
		assert.True(t, ok)
		wg.Done()
	}()

	predefinedSource(src, 1, 2, 3, 4)
	tick()
	wg.Wait()
}
