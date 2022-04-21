package database

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/ezbuy/statsd"
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
