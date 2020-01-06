package strategies

import (
	"github.com/matang28/reshape/reshape"
	"github.com/matang28/reshape/reshape/sinks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirectStrategy_Solve_HappyCase(t *testing.T) {
	strg := NewDirectStrategy()
	src := make(chan interface{})
	sink := sinks.NewArraySink()

	go strg.Solve(src, nil, []interface{}{plusOneTrans, plusOneTrans, dropEvens, sink})

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
	strg := NewDirectStrategy()
	src := make(chan interface{})
	sink := sinks.NewArraySink()
	errors := make(chan error)
	go strg.Solve(src, errors, []interface{}{plusOneTrans, plusOneTrans, badTrans, sink})

	predefinedSource(src, 1, 2, 3, 4)
	tick()

	assert.Nil(t, sink.Get())
	err := <-errors
	err = <-errors
	err = <-errors
	err = <-errors
	_, ok := err.(*reshape.TransformationError)
	assert.True(t, ok)
}

func TestDirectStrategy_Solve_BadSink(t *testing.T) {
	strg := NewDirectStrategy()
	src := make(chan interface{})
	sink := &badSink{}
	errors := make(chan error)

	go strg.Solve(src, errors, []interface{}{plusOneTrans, plusOneTrans, sink})
	predefinedSource(src, 1, 2, 3, 4)
	tick()

	err := <-errors
	err = <-errors
	err = <-errors
	err = <-errors
	_, ok := err.(*reshape.SinkError)
	assert.True(t, ok)
}

func TestDirectStrategy_Solve_UnrecognizedHandler(t *testing.T) {
	strg := NewDirectStrategy()
	src := make(chan interface{})
	errors := make(chan error)

	go strg.Solve(src, errors, []interface{}{plusOneTrans, plusOneTrans, 10})

	predefinedSource(src, 1, 2, 3, 4)
	tick()

	for i := 0; i < 3; i++ {
		err := <-errors
		_, ok := err.(*reshape.UnrecognizedHandler)
		assert.True(t, ok)
	}
}
