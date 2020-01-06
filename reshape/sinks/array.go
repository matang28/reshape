package sinks

type ArraySink struct {
	arr []interface{}
}

func NewArraySink() *ArraySink {
	return &ArraySink{arr: nil}
}

func (this *ArraySink) Dump(object ...interface{}) error {
	this.arr = append(this.arr, object...)
	return nil
}

func (this *ArraySink) CloseGracefully() error {
	return nil
}

func (this *ArraySink) Get() []interface{} {
	return this.arr
}
