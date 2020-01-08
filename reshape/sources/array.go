package sources

import (
	"fmt"
	"github.com/matang28/reshape/reshape"
)

type ArraySource struct {
	array    []interface{}
	ch       chan interface{}
	isClosed bool
}

func NewArraySource() *ArraySource {
	return &ArraySource{ch: make(chan interface{})}
}

func (this *ArraySource) Stream() reshape.Stream {
	return reshape.NewStream(this)
}

func (this *ArraySource) GetChannel() <-chan interface{} {
	return this.ch
}

func (this *ArraySource) Close() error {
	close(this.ch)
	this.isClosed = true
	return nil
}

func (this *ArraySource) Append(elements ...interface{}) error {
	if !this.isClosed {
		this.array = append(this.array, elements...)
		go func() {
			for _, e := range elements {
				this.ch <- e
			}
		}()
		return nil
	}
	return fmt.Errorf("cannot put new entries on closed ArraySource")
}
