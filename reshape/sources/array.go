package sources

import (
	"github.com/matang28/reshape/reshape"
)

type ArraySource struct {
	array []interface{}
	ch    chan interface{}
}

func NewArraySource() *ArraySource {
	return &ArraySource{ch: make(chan interface{})}
}

func (this *ArraySource) Stream() reshape.Stream {
	return reshape.NewStream(this.ch)
}

func (this *ArraySource) CloseGracefully() error {
	close(this.ch)
	return nil
}

func (this *ArraySource) Append(elements ...interface{}) {
	this.array = append(this.array, elements...)
	for _, e := range elements {
		this.ch <- e
	}
}
