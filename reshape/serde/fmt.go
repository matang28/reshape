package serde

import "fmt"

var FmtSerializer Serializer = func(item interface{}) (s string, e error) {
	return fmt.Sprintf("%+v", item), nil
}
