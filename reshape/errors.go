package reshape

import "fmt"

func Report(err error, errors chan error) {
	go func() {
		if err != nil && errors != nil {
			errors <- err
		}
	}()
}

type UnrecognizedHandler struct {
	instance interface{}
}

func NewUnrecognizedHandler(instance interface{}) error {
	return &UnrecognizedHandler{instance: instance}
}

func (this *UnrecognizedHandler) Error() string {
	return fmt.Sprintf("unrecognized handler %#v, currently supported handlers are instances of: Transformation, Filter, Sink. If you did pass one of these make sure that structs are passed as pointers", this.instance)
}

type TransformationError struct {
	E error
}

func NewTransformationError(error error) error {
	if error == nil {
		return nil
	}
	return &TransformationError{E: error}
}

func (t *TransformationError) Error() string {
	if t.E != nil {
		return t.Error()
	}
	return ""
}

type FilterError struct {
	E error
}

func NewFilterError(error error) error {
	if error == nil {
		return nil
	}
	return &FilterError{E: error}
}

func (t *FilterError) Error() string {
	if t.E != nil {
		return t.Error()
	}
	return ""
}

type SinkError struct {
	E error
}

func NewSinkError(error error) error {
	if error == nil {
		return nil
	}
	return &SinkError{E: error}
}

func (t *SinkError) Error() string {
	if t.E != nil {
		return t.Error()
	}
	return ""
}
