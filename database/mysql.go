package database

func newMySQLTracer(options ...TracerOption) *Tracer {
	return NewTracer("mysql", options...)
}

func newMySQLTracerWrapperWithOpts(options ...TracerOption) *TracerWrapper {
	return NewTracerWrapperWithTracer(
		newMySQLTracer(options...),
	)
}

func NewMySQLTracerWrapper() *TracerWrapper {
	return NewTracerWrapper("mysql")
}
