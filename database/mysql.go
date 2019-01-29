package database

func newMySQLTracer(options ...TracerOption) *tracer {
	return newTracer("mysql", options...)
}

// NewMySQLTracerWrapperWithOpts init a pure TracerWrapper with set options
func NewMySQLTracerWrapperWithOpts(options ...TracerOption) *TracerWrapper {
	return newTracerWrapperWithTracer(
		newMySQLTracer(options...),
	)
}

// NewMySQLTracerWrapper init a default TracerWrapper with ignoreSelectColumnsOption
func NewMySQLTracerWrapper() *TracerWrapper {
	return NewTracerWrapper("mysql")
}
