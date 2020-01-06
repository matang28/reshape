package sources

import (
	"github.com/matang28/reshape/reshape/sinks"
	"github.com/matang28/reshape/reshape/strategies"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestArraySource_Stream(t *testing.T) {
	source := NewArraySource()
	sink := sinks.NewArraySink()

	go source.Stream().Filter(filterOdds).Sink(sink).Run(strategies.NewDirectStrategy())
	source.Append(1, 2, 3, 4, 5, 6, "7")

	time.Sleep(100 * time.Millisecond)
	assert.EqualValues(t, []interface{}{2, 4, 6}, sink.Get())
}

var filterOdds = func(in interface{}) bool {
	number, ok := in.(int)
	if !ok {
		return false
	}
	if number%2 == 0 {
		return true
	}
	return false
}
