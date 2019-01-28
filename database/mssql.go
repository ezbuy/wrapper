package database

func newMsSQLTracer(options ...TracerOption) *tracer {
	return newTracer("mssql", options...)
}

// NewMsSQLTracerWrapperWithOpts init a pure TracerWrapper with set options
func NewMsSQLTracerWrapperWithOpts(options ...TracerOption) *TracerWrapper {
	return newTracerWrapperWithTracer(
		newMsSQLTracer(options...),
	)
}

// NewMsSQLTracerWrapper init a default TracerWrapper with ignoreSelectColumnsOption
func NewMsSQLTracerWrapper() *TracerWrapper {
	return NewTracerWrapper("mssql")
}
