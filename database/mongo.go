package database

import (
	"go.mongodb.org/mongo-driver/event"
)

func NewMongoQueryTracer(options ...TracerOption) *TracerWrapper {
	return newTracerWrapperWithTracer(newTracer("mongo", options...))
}

func NewMongoPoolMonitor(t MonitorType) Monitor {
	switch t {
	case StatsD:
		return &StatsDPoolMonitor{
			prefix: "database.mongo",
		}
	case Prometheus:
		return &PrometheusPoolMonitor{
			prefix: "database.mongo",
		}
	default:
		return &DefaultPoolMonitor{}
	}
}

func NewMongoDriverMonitor(m Monitor) *event.PoolMonitor {
	return &event.PoolMonitor{
		Event: func(evt *event.PoolEvent) {
			switch evt.Type {
			case event.PoolCreated:
				m.Pool(PoolCreate)
			case event.PoolCleared:
				m.Pool(PoolClear)
			case event.ConnectionCreated:
				m.Conn(ConnCreate)
			case event.ConnectionClosed:
				m.Conn(ConnClose)
			case event.ConnectionReturned:
				m.Conn(ConnRelease)
			case event.GetSucceeded:
				m.Conn(Connoccupy)
			}
		},
	}
}
