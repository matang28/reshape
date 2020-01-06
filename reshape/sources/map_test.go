package sources

import (
	"github.com/matang28/reshape/reshape/sinks"
	"github.com/matang28/reshape/reshape/strategies"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMapSource_Stream(t *testing.T) {
	expected := map[interface{}]interface{}{1: 11, 2: 22, 3: 33}
	source := NewMapSource()
	sink := sinks.NewMapSink(func(in interface{}) interface{} {
		out, ok := in.(int)
		if !ok {
			return nil
		}
		return out % 10
	})

	go source.Stream().Reshape(func(in interface{}) (out interface{}, err error) {
		o, ok := in.(MapEntry)
		if !ok {
			return nil, nil
		}
		return o.Value, nil
	}).Filter(func(in interface{}) bool {
		return in != nil
	}).Sink(sink).Run(strategies.NewDirectStrategy())

	for k, v := range expected {
		source.Put(k, v)
	}

	time.Sleep(100 * time.Millisecond)
	assert.EqualValues(t, expected, sink.Get())
}
