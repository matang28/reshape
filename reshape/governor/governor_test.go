package governor

import (
	"fmt"
	"github.com/matang28/reshape/reshape"
	"github.com/matang28/reshape/reshape/sinks"
	"github.com/matang28/reshape/reshape/sources"
	"github.com/matang28/reshape/reshape/strategies"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGovernor_StartStop(t *testing.T) {
	// This test will check the start-stop mechanism
	g := New(Config{}, strategies.NewDirectStrategy())

	source := sources.NewArraySource()
	sink := sinks.NewArraySink()
	stream := source.Stream().Reshape(plusOne).Sink(sink)
	g.Add(stream)

	source.Append(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	go func() {
		time.Sleep(1 * time.Second)
		go g.Stop()
	}()

	g.Start()
	assert.EqualValues(t, []interface{}{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, sink.Get())
	source.Append(55) // closed stream should not change anything
	assert.EqualValues(t, []interface{}{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, sink.Get())
}

func TestGovernor_SetErrorHandlerReturnFalse(t *testing.T) {
	// This test will check the error handler, by returning false we are ignoring any kind of error limits:
	g := New(Config{MaxTransformationsErrors: 0}, strategies.NewDirectStrategy())

	source := sources.NewArraySource()
	sink := sinks.NewArraySink()
	stream := source.Stream().Reshape(badTransform).Sink(sink)
	g.Add(stream)
	g.SetErrorHandler(func(err error) bool {
		_, ok := err.(*reshape.TransformationError)
		assert.True(t, ok)
		return false
	})

	source.Append(1, 2, 3, 4)

	go func() {
		time.Sleep(1 * time.Second)
		go g.Stop()
	}()

	g.Start()
	assert.EqualValues(t, []interface{}{1, 3, 4}, sink.Get())
}

func TestGovernor_SetErrorHandlerReturnTrue(t *testing.T) {
	// This test will check the error handler, by returning true the governor should stop the source
	g := New(Config{MaxTransformationsErrors: 0}, strategies.NewDirectStrategy())

	source := sources.NewArraySource()
	sink := sinks.NewArraySink()
	stream := source.Stream().Reshape(badTransform).Sink(sink)
	g.Add(stream)
	g.SetErrorHandler(func(err error) bool {
		_, ok := err.(*reshape.TransformationError)
		assert.True(t, ok)
		return true
	})

	source.Append(1, 2, 3, 4)

	go func() {
		time.Sleep(1 * time.Second)
		go g.Stop()
	}()

	g.Start()
	assert.EqualValues(t, []interface{}{1, 3, 4}, sink.Get())
}

var plusOne = func(in interface{}) (out interface{}, err error) {
	return in.(int) + 1, nil
}

var badTransform = func(in interface{}) (out interface{}, err error) {
	if in.(int) == 2 {
		return nil, fmt.Errorf("bad bad leroy brown")
	}
	return in, nil
}
