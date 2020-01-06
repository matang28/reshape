package sinks

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type person struct {
	Id   string
	Name string
	Age  int
}

func TestMapSink_Dump_ValidKeyExtractor(t *testing.T) {
	a := &person{Id: "a", Name: "Alpha", Age: 1}
	b := &person{Id: "b", Name: "Beta", Age: 2}
	c := &person{Id: "c", Name: "Gamma", Age: 3}
	expectedMap := map[string]interface{}{
		"a": a,
		"b": b,
		"c": c,
	}

	sink := NewMapSink(func(in interface{}) (string, error) {
		return in.(*person).Id, nil
	})

	err := sink.Dump(a, b, c)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedMap, sink.Get())

	bnew := &person{Id: "b", Name: "Beta2", Age: 22}
	d := &person{Id: "d", Name: "Delta", Age: 4}
	expectedMap["b"] = bnew
	expectedMap["d"] = d
	err = sink.Dump(bnew, d)
	assert.Nil(t, err)

	assert.EqualValues(t, expectedMap, sink.Get())
}

func TestMapSink_Dump_InvalidKeyExtractor(t *testing.T) {
	a := &person{Id: "a", Name: "Alpha", Age: 1}
	sink := NewMapSink(func(in interface{}) (string, error) {
		out, ok := in.(string)
		if !ok {
			return "", fmt.Errorf("cannot convert %+v to string", in)
		}
		return out, nil
	})

	err := sink.Dump(a)
	assert.NotNil(t, err)
}
