package sinks

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestFileSink_Dump_NoLimits(t *testing.T) {
	sink := NewFileSink(FileSinkConfig{}, nil)
	person1 := &person{
		Id:   "1",
		Name: "Metuselah",
		Age:  700,
	}
	person2 := &person{
		Id:   "2",
		Name: "Noah",
		Age:  88,
	}
	expected := `{"Id":"1","Name":"Metuselah","Age":700}
{"Id":"2","Name":"Noah","Age":88}
`

	err := sink.Dump(person1, person2)
	assert.Nil(t, err)
	time.Sleep(100 * time.Millisecond)
	bytes, err := ioutil.ReadFile(sink.file.Name())
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	assert.EqualValues(t, expected, string(bytes))

	assert.Nil(t, sink.Close())
	if err := os.Remove(sink.file.Name()); err != nil {
		panic(err)
	}
}
