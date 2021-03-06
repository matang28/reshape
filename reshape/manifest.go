package reshape

// Transformation is a function that takes an object
// and transform it to another object
// failed transformations should return an error
type Transformation = func(in interface{}) (out interface{}, err error)

// Filter is a function that takes an object and returns
// a boolean which indicates if this object should
// be passed to the next handler in the stream.
type Filter = func(in interface{}) bool

// A source is a type that acts as a data source (read)
// for example kafka consumer, CDC systems, files, etc...
type Source interface {
	// Calling Stream() will create a new stream of objects
	// that we can shape, filter and dump to other data sinks.
	Stream() Stream

	// Will return the channel responsible for the communication between the
	// source and the stream
	GetChannel() <-chan interface{}

	// Close will be called by reshape, to signal the data
	// source that it shouldn't generate any new events
	// allowing for the already created events to pass
	// through the pipeline.
	Close() error
}

// A sink is a type that acts as a data source (write)
// for example kafka producer, MySQL, MongoDB, File, etc...
type Sink interface {
	// Will dump these objects into the data sink
	Dump(object ...interface{}) error

	// Close will be called by reshape, to signal the data
	// sink that it shouldn't accept any new events
	// allowing for the already created events to pass
	// through the pipeline.
	Close() error
}

// A stream is a type that acts as a stream of objects.
// Streams can be created from Sources and the objects from that source
// will be processed by the stream pipeline definition.
type Stream interface {
	// Allows you to decide which objects are passed to the next handler
	Filter(filter Filter) Stream

	// Allows you transform the streaming objects to take any other shape.
	Reshape(transformation Transformation) Stream

	// Allows you to stream the transformed objects to custom sink.
	Sink(sink Sink) Stream

	// Will trigger the source to start generating events to be processed by the pipeline.
	Run(strategy StreamingStrategy, errors chan error)

	Close() error
}

// Streaming strategy decides how to process the stream elements
type StreamingStrategy interface {
	// Solve will be called by the stream, it should pull elements from the source channels
	// and process the stream elements with the array of handlers (transformations, filters, sinks, etc...)
	// an error channel is used to report errors thrown while processing the stream.
	Solve(source <-chan interface{}, errors chan error, handlers []interface{})
}
