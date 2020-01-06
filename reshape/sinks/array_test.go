package sinks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArraySink_Dump(t *testing.T) {
	expected := []interface{}{1, 2, "3", 4.5, true}
	sink := NewArraySink()
	err := sink.Dump(expected...)
	assert.Nil(t, err)
	assert.EqualValues(t, expected, sink.Get())

	err = sink.Dump(expected...)
	assert.Nil(t, err)
	assert.EqualValues(t, append(expected, expected...), sink.Get())
}
