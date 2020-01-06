package strategies

import (
	"github.com/matang28/reshape/reshape"
	"github.com/matang28/reshape/reshape/sinks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBufferedStrategy_Solve_CheckBatching(t *testing.T) {
	strg := NewBufferedStrategy(5, 1*time.Minute)
	src := make(chan interface{})
	sink := sinks.NewArraySink()

	go strg.Solve(src, nil, []interface{}{plusOneTrans, plusOneTrans, dropEvens, sink})

	predefinedSource(src, 1, 2, 3, 4)
	tick()
	assert.Empty(t, sink.Get())

	predefinedSource(src, 5, 6, 7, 8)
	tick()
	assert.EqualValues(t, []interface{}{3, 5, 7}, sink.Get())

	predefinedSource(src, 9, 10)
	tick()
	assert.EqualValues(t, []interface{}{3, 5, 7, 9, 11}, sink.Get())
	close(src)
}

func TestBufferedStrategy_Solve_CheckTimeout(t *testing.T) {
	strg := NewBufferedStrategy(5, 50*time.Millisecond)
	src := make(chan interface{})
	sink := sinks.NewArraySink()

	go strg.Solve(src, nil, []interface{}{plusOneTrans, plusOneTrans, dropEvens, sink})

	predefinedSource(src, 1, 2, 3, 4)
	tick()
	assert.EqualValues(t, []interface{}{3, 5}, sink.Get())

	predefinedSource(src, 5, 6, 7, 8)
	tick()
	assert.EqualValues(t, []interface{}{3, 5, 7, 9}, sink.Get())

	predefinedSource(src, 9)
	tick()
	assert.EqualValues(t, []interface{}{3, 5, 7, 9, 11}, sink.Get())
	close(src)
}

func TestBufferedStrategy_Solve_BadTransformation(t *testing.T) {
	strg := NewBufferedStrategy(5, 50*time.Millisecond)
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

func TestBufferedStrategy_Solve_BadSink(t *testing.T) {
	strg := NewBufferedStrategy(5, 50*time.Millisecond)
	src := make(chan interface{})
	sink := &badSink{}
	errors := make(chan error)

	go strg.Solve(src, errors, []interface{}{plusOneTrans, plusOneTrans, sink})
	predefinedSource(src, 1, 2, 3, 4)
	tick()

	err := <-errors
	_, ok := err.(*reshape.SinkError)
	assert.True(t, ok)
}

func TestBufferedStrategy_Solve_UnrecognizedHandler(t *testing.T) {
	strg := NewBufferedStrategy(5, 50*time.Millisecond)
	src := make(chan interface{})
	errors := make(chan error)

	go strg.Solve(src, errors, []interface{}{plusOneTrans, plusOneTrans, 10})

	predefinedSource(src, 1, 2, 3, 4)
	tick()

	err := <-errors
	_, ok := err.(*reshape.UnrecognizedHandler)
	assert.True(t, ok)
}
