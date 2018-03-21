package main

import (
	"time"

	"github.com/sirupsen/logrus"
)

type metrics map[string]*metricsEntry

type metricsEntry struct {
	start   time.Time
	elapsed time.Duration
}

func newMetrics() *metrics {
	return &metrics{}
}

func (m metrics) newEntry(name string) *metricsEntry {
	entry := &metricsEntry{
		start: time.Now(),
	}

	m[name] = entry

	return entry
}

func (m metrics) log() {
	val := make(map[string]interface{}, len(m))
	for k, v := range m {
		val[k] = v.elapsed
	}

	logrus.WithFields(val).Debug("metrics")
}

func (m *metricsEntry) done() {
	m.elapsed = time.Since(m.start)
}
