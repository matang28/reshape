package reshape

import "fmt"

type UnrecognizedHandler struct {
	instance interface{}
}

func NewUnrecognizedHandler(instance interface{}) error {
	return &UnrecognizedHandler{instance: instance}
}

func (this *UnrecognizedHandler) Error() string {
	return fmt.Sprintf("unrecognized handler %#v, currently supported handlers are instances of: Transformation, Filter, Sink. If you did pass one of these make sure that structs are passed as pointers", this.instance)
}
