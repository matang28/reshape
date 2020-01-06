package sinks

type KeyExtractor = func(in interface{}) (string, error)

type MapSink struct {
	keyExtractor KeyExtractor
	m            map[string]interface{}
}

func NewMapSink(keyExtractor KeyExtractor) *MapSink {
	return &MapSink{keyExtractor: keyExtractor, m: make(map[string]interface{})}
}

func (this *MapSink) Dump(object ...interface{}) error {
	for _, elem := range object {
		key, err := this.keyExtractor(elem)
		if err != nil {
			return err
		}
		this.m[key] = elem
	}
	return nil
}

func (this *MapSink) CloseGracefully() error {
	return nil
}

func (this *MapSink) Get() map[string]interface{} {
	return this.m
}
