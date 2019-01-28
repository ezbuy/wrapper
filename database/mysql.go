package database

func newMySQLTracer(options ...WrapperOption) *Tracer {
	return NewTracer("mysql", options...)
}

func newMySQLTracerWrapperWithOpts(options ...WrapperOption) *TracerWrapper {
	return NewTracerWrapperWithTracer(
		newMySQLTracer(options...),
	)
}

func NewMySQLTracerWrapper() *TracerWrapper {
	return NewTracerWrapper("mysql")
}
