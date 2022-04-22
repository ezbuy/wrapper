package database

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/ezbuy/statsd"
	"github.com/ezbuy/wrapper/pkg/net"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
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

type Monitor interface {
	Pool(PoolOperation) error
	Conn(ConnOperation) error
}

type StatsDPoolMonitor struct {
	prefix string
}

func (c *StatsDPoolMonitor) Pool(op PoolOperation) error {
	switch op {
	case PoolCreate:
		statsd.Incr(c.prefix + ".pool")
	case PoolClear:
		statsd.IncrByVal(c.prefix+".pool", -1)
	}
	return nil
}

func (c *StatsDPoolMonitor) Conn(op ConnOperation) error {
	switch op {
	case ConnCreate:
		statsd.Incr(c.prefix + ".conn")
	case ConnClose:
		statsd.IncrByVal(c.prefix+".conn", -1)
	case Connoccupy:
		statsd.Incr(c.prefix + ".conn.occupy")
	case ConnRelease:
		statsd.IncrByVal(c.prefix+".conn.occupy", -1)
	}
	return nil
}

// PrometheusMonitor is the prometheus based monitor
// but must be used with prometheus pushgateway
type PrometheusPoolMonitor struct {
	prefix string
	p      *push.Pusher
	reg    *prometheus.Registry
}

const (
	subsystemScope = "m_database_pool"
)

var (
	poolCreate = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: subsystemScope,
		Name:      "pool_create",
		Help:      "pool create",
	})

	poolClear = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: subsystemScope,
		Name:      "mongo_pool_clear",
		Help:      "pool clear",
	})
	conn = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: subsystemScope,
		Name:      "mongo_conn",
		Help:      "conn create",
	})
	connOccupy = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: subsystemScope,
		Name:      "mongo_conn_occupy",
		Help:      "conn occupy",
	})
)

func NewPrometheusPoolMonitor(gatewayAddress string) *PrometheusPoolMonitor {
	reg := prometheus.NewRegistry()
	reg.MustRegister(poolCreate, poolClear, conn, connOccupy)

	return &PrometheusPoolMonitor{
		prefix: "database.mongo",
		p:      push.New(gatewayAddress, "database.mongo"),
		reg:    reg,
	}
}

func (c *PrometheusPoolMonitor) push() error {
	if err := c.p.Gatherer(c.reg).Grouping(
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
		poolCreate.Inc()
	case PoolClear:
		poolClear.Dec()
	}
	return c.push()
}

func (c *PrometheusPoolMonitor) Conn(op ConnOperation) error {
	switch op {
	case ConnCreate:
		conn.Inc()
	case ConnClose:
		conn.Dec()
	case Connoccupy:
		connOccupy.Inc()
	case ConnRelease:
		connOccupy.Dec()
	}
	return c.push()
}

type DefaultPoolMonitor struct {
	poolSize  int64
	connNum   int64
	occupyNum int64
	sync.Mutex
}

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
