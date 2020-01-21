package resources

import (
	"errors"
	"fmt"

	prometheus "github.com/prometheus/client_golang/prometheus"
	customMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

// Set exports a Prometheus Gauge
func SetMetric(name string, labels map[string]string, epochTimestamp int) error {

	metric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        name,
			Help:        labels["description"],
			ConstLabels: labels,
		},
	)

	if err := customMetrics.Registry.Unregister(metric); err != false {
		return errors.New(fmt.Sprintf("Unable to unregister Prometheus metric: %s", metric))
	}
	customMetrics.Registry.MustRegister(metric)

	return nil
}
