package database

func newMsSQLTracer(options ...func(*Tracer)) *Tracer {
	return NewCustmizedTracer("mssql", options...)
}

func NewMsSQLTracerWrapper(options ...func(*Tracer)) *TracerWrapper {
	return NewCustmizedTracerWrapper(
		newMsSQLTracer(options...),
	)
}

func NewDefaultMsSQLTracerWrapper() *TracerWrapper {
	return NewTracerWrapper("mssql")
}
