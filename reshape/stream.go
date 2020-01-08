package reshape

type stream struct {
	handlers   []interface{}
	sourceChan <-chan interface{}
	source     Source
}

func NewStream(source Source) Stream {
	return &stream{sourceChan: source.GetChannel(), source: source}
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

func (this *stream) Run(strategy StreamingStrategy, errors chan error) {
	go strategy.Solve(this.sourceChan, errors, this.handlers)
}

func (this *stream) Close() error {
	return this.source.Close()
}
