package governor

// Governor can expose its stats thorough this struct
type Stats struct {
	TransformationsErrors int
	FilterErrors          int
	SinkErrors            int
}
