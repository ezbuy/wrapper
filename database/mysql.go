package database

func newMySQLTracer(options ...func(*Tracer)) *Tracer {
	return NewCustmizedTracer("mysql", options...)
}

func NewMySQLTracerWrapper(options ...func(*Tracer)) *TracerWrapper {
	return NewCustmizedTracerWrapper(
		newMySQLTracer(options...),
	)
}

func NewDefaultMySQLTracerWrapper() *TracerWrapper {
	return NewTracerWrapper("mysql")
}
