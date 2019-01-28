package database

func newMsSQLTracer(options ...TracerOption) *Tracer {
	return NewTracer("mssql", options...)
}

func newMsSQLTracerWrapperWithOpts(options ...TracerOption) *TracerWrapper {
	return NewTracerWrapperWithTracer(
		newMsSQLTracer(options...),
	)
}

func NewMsSQLTracerWrapper() *TracerWrapper {
	return NewTracerWrapper("mssql")
}
