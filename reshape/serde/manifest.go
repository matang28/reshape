package serde

type Serializer func(item interface{}) (string, error)

type Deserializer func(str string, out interface{}) error
