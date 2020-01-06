package reshape

type stream struct {
	handlers   []interface{}
	sourceChan <-chan interface{}
}

func NewStream(sourceChan <-chan interface{}) Stream {
	return &stream{sourceChan: sourceChan}
}

func (this *stream) Filter(filter Filter) Stream {
	this.handlers = append(this.handlers, filter)
	return this
}

func (this *stream) Reshape(transformation Transformation) Stream {
	this.handlers = append(this.handlers, transformation)
	return this
}

func (this *stream) Sink(sink Sink) Stream {
	this.handlers = append(this.handlers, sink)
	return this
}

func (this *stream) Run(strategy StreamingStrategy) error {
	return strategy.Solve(this.sourceChan, this.handlers)
}
