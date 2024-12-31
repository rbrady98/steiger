package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	registry *prometheus.Registry

	JokesCreated prometheus.Counter
}

func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry()

	m := &Metrics{
		registry: reg,
		JokesCreated: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "steiger",
			Name:      "jokes_created",
		}),
	}

	reg.MustRegister(m.JokesCreated)
	reg.MustRegister(collectors.NewGoCollector())

	return m
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
