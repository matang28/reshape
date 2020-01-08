package governor

import (
	"fmt"
	"github.com/matang28/reshape/reshape"
	"log"
	"sync"
)

// Governor is responsible for running multiple streams
// giving them a monitored environment to run
// although streams can be managed manually you will have to handle the go-routine synchronisation, error handling and monitoring
// governors abstract this complexity away by giving you a simple api for handling streams
type governor struct {
	streams  []reshape.Stream
	strategy reshape.StreamingStrategy

	config Config
	stats  Stats

	errors       chan error
	errorHandler ErrorHandler

	wg sync.WaitGroup
}

func New(config Config, strategy reshape.StreamingStrategy) *governor {
	return &governor{
		strategy:     strategy,
		errors:       make(chan error),
		errorHandler: defaultErrorHandler,
		stats:        Stats{},
		config:       config,
		wg:           sync.WaitGroup{},
	}
}

// Add stream to this governor
func (this *governor) Add(stream reshape.Stream) {
	this.streams = append(this.streams, stream)
}

// Sets the error handler that the governor uses (check ErrorHandler for more info)
func (this *governor) SetErrorHandler(handler ErrorHandler) {
	this.errorHandler = handler
}

// Starts processing all streams blocking until the stream has stopped
// either by explicit call to governor.Stop() or when error limit is reached
func (this *governor) Start() {
	for _, s := range this.streams {
		s.Run(this.strategy, this.errors)
		this.wg.Add(1)
	}
	go this.handleErrorLimits()
	this.wg.Wait()
}

// Stops the processing of all streams
func (this *governor) Stop() {
	defer func() {
		if d := recover(); d != nil {
			log.Println(d)
		}
	}()

	for _, s := range this.streams {
		if err := s.Close(); err != nil {
			log.Println(err)
		}
		this.wg.Done()
	}
}

// Get the copy of the current stats:
func (this *governor) GetStats() Stats {
	return this.stats
}

func (this *governor) handleErrorLimits() {
	for {
		err := <-this.errors
		if this.errorHandler(err) {
			switch err.(type) {
			case *reshape.TransformationError:
				this.stats.TransformationsErrors++
			case *reshape.FilterError:
				this.stats.FilterErrors++
			case *reshape.SinkError:
				this.stats.SinkErrors++
			}
		}

		if this.config.MaxSinkErrors >= 0 && this.stats.SinkErrors > this.config.MaxSinkErrors {
			log.Println(fmt.Sprintf("Max sink errors has reached (%d)! Stopping reshape...", this.config.MaxSinkErrors))
			this.Stop()
		}

		if this.config.MaxTransformationsErrors >= 0 && this.stats.TransformationsErrors > this.config.MaxTransformationsErrors {
			log.Println(fmt.Sprintf("Max transformation errors has reached (%d)! Stopping reshape...", this.config.MaxTransformationsErrors))
			this.Stop()
		}

		if this.config.MaxFilterErrors >= 0 && this.stats.FilterErrors > this.config.MaxFilterErrors {
			log.Println(fmt.Sprintf("Max filter errors has reached (%d)! Stopping reshape...", this.config.MaxFilterErrors))
			this.Stop()
		}
	}
}
