package database

import (
	"fmt"
	"io"
	"log"
	"sync"
	"sync/atomic"

	"github.com/ezbuy/statsd"
	"github.com/ezbuy/wrapper/pkg/net"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const (
	DEBUG_ENV = "DEBUG_MONITOR"
)

type MonitorType uint8

const (
	StatsD MonitorType = iota + 1
	Prometheus
)

type ConnOperation uint8

const (
	ConnCreate ConnOperation = iota + 1
	ConnClose
	Connoccupy
	ConnRelease
)

type PoolOperation uint8

const (
	PoolCreate PoolOperation = iota + 1
	PoolClear
)

type Logger interface {
	Log(io.Writer, interface{})
}

type Monitor interface {
	Logger
	Pool(PoolOperation) error
	Conn(ConnOperation) error
}

type StatsDPoolMonitor struct {
	prefix string
}

func NewStatsDPoolMonitor(appName string) *StatsDPoolMonitor {
	return &StatsDPoolMonitor{
		prefix: appName,
	}
}

func (m *StatsDPoolMonitor) Log(w io.Writer, args interface{}) {
	fmt.Fprintf(w, "statsd pool monitor: %v", args)
}

func (c *StatsDPoolMonitor) Pool(op PoolOperation) error {
	switch op {
	case PoolCreate:
		statsd.Incr(c.prefix + ".db.pool")
	case PoolClear:
		statsd.IncrByVal(c.prefix+".db.pool", -1)
	}
	return nil
}

func (c *StatsDPoolMonitor) Conn(op ConnOperation) error {
	switch op {
	case ConnCreate:
		statsd.Incr(c.prefix + ".db.conn")
	case ConnClose:
		statsd.IncrByVal(c.prefix+".db.conn", -1)
	case Connoccupy:
		statsd.Incr(c.prefix + ".db.conn.occupy")
	case ConnRelease:
		statsd.IncrByVal(c.prefix+".db.conn.occupy", -1)
	}
	return nil
}

// PrometheusMonitor is the prometheus based monitor
// but must be used with prometheus pushgateway
type PrometheusPoolMonitor struct {
	prefix     string
	p          *push.Pusher
	collectors map[monitorKind]prometheus.Collector
	sync.Mutex
}

const (
	subsystemScope = "monitor_mongo"
)

type monitorKind uint8

const (
	monitorKindPool monitorKind = iota + 1
	monitorKindConn
	monitorKindConnOccupy
)

func newMonitorPool(prefix string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: subsystemScope,
		Name:      prefix + "_pool_current",
		Help:      "pool ",
	})
}

func newMonitorConn(prefix string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: subsystemScope,
		Name:      prefix + "_conn_current",
		Help:      "conn",
	})
}

func newMonitorConnOccupy(prefix string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: subsystemScope,
		Name:      prefix + "_conn_occupy_current",
		Help:      "conn occupy",
	})
}

func NewPrometheusPoolMonitor(appName string, gatewayAddress string) *PrometheusPoolMonitor {
	reg := prometheus.NewRegistry()
	pool, conn, connOccupy := newMonitorPool(appName), newMonitorConn(appName), newMonitorConnOccupy(appName)
	reg.MustRegister(
		pool, conn, connOccupy,
	)

	return &PrometheusPoolMonitor{
		prefix: appName,
		p:      push.New(gatewayAddress, "mongo-pool-monitor").Gatherer(reg),
		collectors: map[monitorKind]prometheus.Collector{
			monitorKindPool:       pool,
			monitorKindConn:       conn,
			monitorKindConnOccupy: connOccupy,
		},
	}
}

func (m *PrometheusPoolMonitor) Log(w io.Writer, args interface{}) {
	fmt.Fprintf(w, "prometheus pool monitor: %v", args)
}

func (c *PrometheusPoolMonitor) push() error {
	if err := c.p.Grouping(
		"kind", "mongo",
	).Grouping(
		"instance", net.GetOutboundIP(),
	).Push(); err != nil {
		return err
	}
	return nil
}

func (c *PrometheusPoolMonitor) Pool(op PoolOperation) error {
	switch op {
	case PoolCreate:
		c.collectors[monitorKindPool].(prometheus.Gauge).Inc()
	case PoolClear:
		c.collectors[monitorKindPool].(prometheus.Gauge).Dec()
	}
	c.Lock()
	defer c.Unlock()
	return c.push()
}

func (c *PrometheusPoolMonitor) Conn(op ConnOperation) error {
	switch op {
	case ConnCreate:
		c.collectors[monitorKindConn].(prometheus.Gauge).Inc()
	case ConnClose:
		c.collectors[monitorKindConn].(prometheus.Gauge).Dec()
	case Connoccupy:
		c.collectors[monitorKindConnOccupy].(prometheus.Gauge).Inc()
	case ConnRelease:
		c.collectors[monitorKindConnOccupy].(prometheus.Gauge).Dec()
	}
	c.Lock()
	defer c.Unlock()
	return c.push()
}

type DefaultPoolMonitor struct {
	poolSize  int64
	connNum   int64
	occupyNum int64
	sync.Mutex
}

func (c *DefaultPoolMonitor) Log(w io.Writer, a interface{}) {}

func (c *DefaultPoolMonitor) Pool(op PoolOperation) error {
	var s int64
	switch op {
	case PoolCreate:
		s = atomic.AddInt64(&c.poolSize, 1)
	case PoolClear:
		s = atomic.AddInt64(&c.poolSize, -1)
	}
	log.Printf("current pool size: %d", s)
	c.Lock()
	defer c.Unlock()
	c.poolSize = s
	return nil
}

func (c *DefaultPoolMonitor) Conn(op ConnOperation) error {
	var s, o int64
	switch op {
	case ConnCreate:
		s = atomic.AddInt64(&c.connNum, 1)
	case ConnClose:
		s = atomic.AddInt64(&c.connNum, -1)
	case Connoccupy:
		o = atomic.AddInt64(&c.occupyNum, 1)
	case ConnRelease:
		o = atomic.AddInt64(&c.occupyNum, -1)
	}
	log.Printf("current pool size: %d, occupy : %d", s, o)
	c.Lock()
	defer c.Unlock()
	c.connNum = s
	c.occupyNum = o
	return nil
}
