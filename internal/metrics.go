package internal

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

// DefObjectives ...
var DefObjectives = map[float64]float64{
	0.50: 0.05,
	0.90: 0.01,
	0.95: 0.005,
	0.99: 0.001,
}

// Metrics ...
type Metrics struct {
	sync.RWMutex

	namespace string
	metrics   map[string]prometheus.Metric
	guagevecs map[string]*prometheus.GaugeVec
	sumvecs   map[string]*prometheus.SummaryVec
}

// NewMetrics ...
func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		namespace: namespace,
		metrics:   make(map[string]prometheus.Metric),
		guagevecs: make(map[string]*prometheus.GaugeVec),
		sumvecs:   make(map[string]*prometheus.SummaryVec),
	}
}

// NewCounter ...
func (m *Metrics) NewCounter(subsystem, name, help string) prometheus.Counter {
	counter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: m.namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
	)

	key := fmt.Sprintf("%s_%s", subsystem, name)
	m.Lock()
	m.metrics[key] = counter
	m.Unlock()
	prometheus.MustRegister(counter)

	return counter
}

// NewCounterFunc ...
func (m *Metrics) NewCounterFunc(subsystem, name, help string, f func() float64) prometheus.CounterFunc {
	counter := prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Namespace: m.namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
		f,
	)

	key := fmt.Sprintf("%s_%s", subsystem, name)
	m.Lock()
	m.metrics[key] = counter
	m.Unlock()
	prometheus.MustRegister(counter)

	return counter
}

// NewGauge ...
func (m *Metrics) NewGauge(subsystem, name, help string) prometheus.Gauge {
	guage := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
	)

	key := fmt.Sprintf("%s_%s", subsystem, name)
	m.Lock()
	m.metrics[key] = guage
	m.Unlock()
	prometheus.MustRegister(guage)

	return guage
}

// NewGaugeFunc ...
func (m *Metrics) NewGaugeFunc(subsystem, name, help string, f func() float64) prometheus.GaugeFunc {
	guage := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
		f,
	)

	key := fmt.Sprintf("%s_%s", subsystem, name)
	m.Lock()
	m.metrics[key] = guage
	m.Unlock()
	prometheus.MustRegister(guage)

	return guage
}

// NewGaugeVec ...
func (m *Metrics) NewGaugeVec(subsystem, name, help string, labels []string) *prometheus.GaugeVec {
	guagevec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: m.namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
		labels,
	)

	key := fmt.Sprintf("%s_%s", subsystem, name)
	m.Lock()
	m.guagevecs[key] = guagevec
	m.Unlock()
	prometheus.MustRegister(guagevec)

	return guagevec
}

// NewSummary ...
func (m *Metrics) NewSummary(subsystem, name, help string) prometheus.Summary {
	summary := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace:  m.namespace,
			Subsystem:  subsystem,
			Name:       name,
			Help:       help,
			Objectives: DefObjectives,
		},
	)

	key := fmt.Sprintf("%s_%s", subsystem, name)
	m.Lock()
	m.metrics[key] = summary
	m.Unlock()
	prometheus.MustRegister(summary)

	return summary
}

// NewSummaryVec ...
func (m *Metrics) NewSummaryVec(subsystem, name, help string, labels []string) *prometheus.SummaryVec {
	sumvec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  m.namespace,
			Subsystem:  subsystem,
			Name:       name,
			Help:       help,
			Objectives: DefObjectives,
		},
		labels,
	)

	key := fmt.Sprintf("%s_%s", subsystem, name)
	m.Lock()
	m.sumvecs[key] = sumvec
	m.Unlock()
	prometheus.MustRegister(sumvec)

	return sumvec
}

// Counter ...
func (m *Metrics) Counter(subsystem, name string) prometheus.Counter {
	key := fmt.Sprintf("%s_%s", subsystem, name)
	return m.metrics[key].(prometheus.Counter)
}

// Gauge ...
func (m *Metrics) Gauge(subsystem, name string) prometheus.Gauge {
	key := fmt.Sprintf("%s_%s", subsystem, name)
	m.RLock()
	defer m.RUnlock()
	return m.metrics[key].(prometheus.Gauge)
}

// GaugeVec ...
func (m *Metrics) GaugeVec(subsystem, name string) *prometheus.GaugeVec {
	key := fmt.Sprintf("%s_%s", subsystem, name)
	m.RLock()
	defer m.RUnlock()
	return m.guagevecs[key]
}

// Summary ...
func (m *Metrics) Summary(subsystem, name string) prometheus.Summary {
	key := fmt.Sprintf("%s_%s", subsystem, name)
	m.RLock()
	defer m.RUnlock()
	return m.metrics[key].(prometheus.Summary)
}

// SummaryVec ...
func (m *Metrics) SummaryVec(subsystem, name string) *prometheus.SummaryVec {
	key := fmt.Sprintf("%s_%s", subsystem, name)
	m.RLock()
	defer m.RUnlock()
	return m.sumvecs[key]
}

// Handler ...
func (m *Metrics) Handler() http.Handler {
	return promhttp.Handler()
}

// Run ...
func (m *Metrics) Run(addr string) {
	http.Handle("/", m.Handler())
	log.Infof("metrics endpoint listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
