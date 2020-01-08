package strategies

import "github.com/matang28/reshape/reshape"

type directStrategy struct {
}

func NewDirectStrategy() *directStrategy {
	return &directStrategy{}
}

func (this *directStrategy) Solve(source <-chan interface{}, errors chan error, handlers []interface{}) {
SourceLoop:
	for {
		var temp, ok = <-source
		if !ok {
			return
		}

		for _, handler := range handlers {
			if temp == nil {
				continue
			}

			switch handler.(type) {
			case reshape.Transformation:
				res := this.transform(handler, temp, errors)
				temp = res

			case reshape.Filter:
				ok := this.filter(handler, temp, errors)
				if !ok {
					continue SourceLoop
				}

			case reshape.Sink:
				this.sink(handler, temp, errors)

			default:
				reshape.Report(reshape.NewUnrecognizedHandler(handler), errors)
			}
		}
	}
}

func (this *directStrategy) filter(handler interface{}, input interface{}, errors chan error) bool {
	defer func() {
		if p := recover(); p != nil {
			err, ok := p.(error)
			if ok {
				reshape.Report(reshape.NewFilterError(err), errors)
			}
		}
	}()

	ok := handler.(reshape.Filter)(input)
	return ok
}

func (this *directStrategy) sink(handler interface{}, input interface{}, errors chan error) {
	defer func() {
		if p := recover(); p != nil {
			err, ok := p.(error)
			if ok {
				reshape.Report(reshape.NewSinkError(err), errors)
			}
		}
	}()

	err := handler.(reshape.Sink).Dump(input)
	reshape.Report(reshape.NewSinkError(err), errors)
}

func (this *directStrategy) transform(handler interface{}, input interface{}, errors chan error) interface{} {
	defer func() {
		if p := recover(); p != nil {
			err, ok := p.(error)
			if ok {
				reshape.Report(reshape.NewTransformationError(err), errors)
			}
		}
	}()

	x, err := handler.(reshape.Transformation)(input)
	reshape.Report(reshape.NewTransformationError(err), errors)
	return x
}
