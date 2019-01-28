package database

func newMsSQLTracer(options ...WrapperOption) *Tracer {
	return NewTracer("mssql", options...)
}

func newMsSQLTracerWrapperWithOpts(options ...WrapperOption) *TracerWrapper {
	return NewTracerWrapperWithTracer(
		newMsSQLTracer(options...),
	)
}

func NewMsSQLTracerWrapper() *TracerWrapper {
	return NewTracerWrapper("mssql")
}
