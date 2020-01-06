package strategies

import (
	"github.com/matang28/reshape/reshape"
	"sync"
	"time"
)

type bufferedStrategy struct {
	batchSize    int
	flushTimeout time.Duration
	queue        chan []interface{}
	mutex        sync.Mutex
}

func NewBufferedStrategy(batchSize int, flushTimeout time.Duration) *bufferedStrategy {
	return &bufferedStrategy{batchSize: batchSize, flushTimeout: flushTimeout, queue: make(chan []interface{}), mutex: sync.Mutex{}}
}

func (this *bufferedStrategy) Solve(source <-chan interface{}, errors chan error, handlers []interface{}) {
	var batch = make([]interface{}, 0)
	go processBatches(this.queue, errors, handlers)

	for {
		select {
		case item := <-source:
			batch = append(batch, item)
			if len(batch) == this.batchSize {
				this.flush(&batch, errors)
			}

		case <-time.After(this.flushTimeout):
			this.flush(&batch, errors)
		}
	}
}

func (this *bufferedStrategy) flush(batch *[]interface{}, errors chan error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if batch != nil && len(*batch) > 0 {
		this.queue <- *batch
		*batch = make([]interface{}, 0)
	}
}

func processBatches(queue chan []interface{}, errors chan error, handlers []interface{}) {
	defer func() {
		d := recover()
		if d != nil {
			e, ok := d.(error)
			if ok {
				reshape.Report(e, errors)
			}
		}
	}()

	for {
		batch := <-queue
		if batch == nil {
			continue
		}

		for _, handler := range handlers {
			switch handler.(type) {
			case reshape.Transformation:
				runTransform(handler.(reshape.Transformation), &batch, errors)
			case reshape.Filter:
				runFilter(handler.(reshape.Filter), &batch, errors)
			case reshape.Sink:
				runSink(handler, batch, errors)
			default:
				reshape.Report(reshape.NewUnrecognizedHandler(handler), errors)
			}
		}
	}
}

func runSink(handler interface{}, batch []interface{}, errors chan error) {
	defer func() {
		if p := recover(); p != nil {
			err, ok := p.(error)
			if ok {
				reshape.Report(reshape.NewSinkError(err), errors)
			}
		}
	}()

	var dump []interface{}
	for _, elem := range batch {
		if elem != nil {
			dump = append(dump, elem)
		}
	}

	err := handler.(reshape.Sink).Dump(dump...)
	reshape.Report(reshape.NewSinkError(err), errors)
}

func runTransform(transformation reshape.Transformation, batch *[]interface{}, errors chan error) {
	defer func() {
		if p := recover(); p != nil {
			err, ok := p.(error)
			if ok {
				reshape.Report(reshape.NewTransformationError(err), errors)
			}
		}
	}()

	for i := 0; i < len(*batch); i++ {
		if (*batch)[i] == nil {
			continue
		}
		out, err := transformation((*batch)[i])
		reshape.Report(reshape.NewTransformationError(err), errors)
		(*batch)[i] = out
	}
}

func runFilter(filter reshape.Filter, batch *[]interface{}, errors chan error) {
	defer func() {
		if p := recover(); p != nil {
			err, ok := p.(error)
			if ok {
				reshape.Report(reshape.NewFilterError(err), errors)
			}
		}
	}()

	var out []interface{}
	for i := 0; i < len(*batch); i++ {
		ok := filter((*batch)[i])
		if ok {
			out = append(out, (*batch)[i])
		}
	}
	*batch = out
}
