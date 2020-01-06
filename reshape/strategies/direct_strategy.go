package strategies

import "github.com/matang28/reshape/reshape"

type directStrategy struct {
}

func NewDirectStrategy() *directStrategy {
	return &directStrategy{}
}

func (this *directStrategy) Solve(source <-chan interface{}, handlers []interface{}) error {
SourceLoop:
	for item := range source {
		var temp = item
		for _, handler := range handlers {
			switch handler.(type) {
			case reshape.Transformation:
				x, err := handler.(reshape.Transformation)(temp)
				if err != nil {
					return err
				}
				temp = x

			case reshape.Filter:
				ok := handler.(reshape.Filter)(temp)
				if !ok {
					continue SourceLoop
				}

			case reshape.Sink:
				if err := handler.(reshape.Sink).Dump(temp); err != nil {
					return err
				}
			default:
				return reshape.NewUnrecognizedHandler(handler)
			}
		}
	}
	return nil
}
