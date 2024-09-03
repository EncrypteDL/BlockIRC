package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	assert := assert.New(t)

	m := NewMetrics("test")
	m.NewCounter("foo", "counter", "help")
	m.NewCounterFunc("foo", "counter_func", "help", func() float64 { return 1.0 })
	m.NewGauge("foo", "gauge", "help")
	m.NewGaugeFunc("foo", "gauge_func", "help", func() float64 { return 1.0 })
	m.NewGaugeVec("foo", "gauge_vec", "help", []string{"test"})

	m.Counter("foo", "counter").Inc()
	m.Gauge("foo", "gauge").Add(1)
	m.GaugeVec("foo", "gauge_vec").WithLabelValues("test").Add(1)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)

	m.Handler().ServeHTTP(w, r)
	assert.Equal(w.Code, http.StatusOK)

	assert.Regexp(
		`
# HELP test_foo_counter help
# TYPE test_foo_counter counter
test_foo_counter 1
# HELP test_foo_counter_func help
# TYPE test_foo_counter_func counter
test_foo_counter_func 1
# HELP test_foo_gauge help
# TYPE test_foo_gauge gauge
test_foo_gauge 1
# HELP test_foo_gauge_func help
# TYPE test_foo_gauge_func gauge
test_foo_gauge_func 1
# HELP test_foo_gauge_vec help
# TYPE test_foo_gauge_vec gauge
test_foo_gauge_vec{test="test"} 1
`,
		w.Body.String(),
	)
}
