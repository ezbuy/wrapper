package database

import (
	"github.com/ezbuy/wrapper"
)

func newMsSQLTracer(options ...func(*Tracer)) *Tracer {
	return NewCustmizedTracer("mssql", options...)
}

func NewMsSQLTracerWrapper(options ...func(*Tracer)) wrapper.Wrapper {
	return NewCustmizedTracerWrapper(
		newMsSQLTracer(options...),
	)
}

func NewDefaultMsSQLTracerWrapper() wrapper.Wrapper {
	return NewDefaultTracerWrapper("mssql")
}
