package strategies

import (
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
