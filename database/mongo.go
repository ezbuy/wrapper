package database

import (
	"fmt"
	"os"

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
			if os.Getenv(DEBUG_ENV) != "" {
				m.Log(os.Stdout, fmt.Sprintf("event: %v\n", evt))
			}
			switch evt.Type {
			case event.PoolCreated:
				if err := m.Pool(PoolCreate); err != nil && os.Getenv(DEBUG_ENV) != "" {
					m.Log(os.Stderr, fmt.Sprintf("pool: %s\n", err))
				}
			case event.PoolCleared:
				if err := m.Pool(PoolClear); err != nil && os.Getenv(DEBUG_ENV) != "" {
					m.Log(os.Stderr, fmt.Sprintf("pool: %s\n", err))
				}
			case event.ConnectionCreated:
				if err := m.Conn(ConnCreate); err != nil && os.Getenv(DEBUG_ENV) != "" {
					m.Log(os.Stderr, fmt.Sprintf("conn: %s\n", err))
				}
			case event.ConnectionClosed:
				if err := m.Conn(ConnClose); err != nil && os.Getenv(DEBUG_ENV) != "" {
					m.Log(os.Stderr, fmt.Sprintf("conn: %s\n", err))
				}
			case event.ConnectionReturned:
				if err := m.Conn(ConnRelease); err != nil && os.Getenv(DEBUG_ENV) != "" {
					m.Log(os.Stderr, fmt.Sprintf("conn: %s\n", err))
				}
			case event.GetSucceeded:
				if err := m.Conn(Connoccupy); err != nil && os.Getenv(DEBUG_ENV) != "" {
					m.Log(os.Stderr, fmt.Sprintf("conn: %s\n", err))
				}
			}
		},
	}
}
