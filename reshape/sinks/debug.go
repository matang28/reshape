package sinks

import (
	"fmt"
	"github.com/matang28/reshape/reshape/serde"
)

type DebugSink struct {
	serializer serde.Serializer
}

func NewDebugSink() *DebugSink {
	return &DebugSink{serializer: serde.FmtSerializer}
}

func NewCustomDebugSink(serializer serde.Serializer) *DebugSink {
	return &DebugSink{serializer: serializer}
}

func (this *DebugSink) Dump(objects ...interface{}) error {
	for _, o := range objects {
		str, err := this.serializer(o)
		if err != nil {
			return err
		}
		fmt.Printf(str + "\n")
	}
	return nil
}

func (this *DebugSink) Close() error {
	return nil
}
