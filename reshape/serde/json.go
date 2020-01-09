package serde

import "encoding/json"

var JsonSerializer Serializer = func(item interface{}) (s string, e error) {
	bytes, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

var JsonDeserializer Deserializer = func(str string, out interface{}) error {
	err := json.Unmarshal([]byte(str), out)
	if err != nil {
		return err
	}
	return nil
}
