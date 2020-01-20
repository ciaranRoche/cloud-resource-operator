package resources

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	customMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	sleepytime = 3600
)

var (
	// create the map of vectors
	MetricVecs map[string]prometheus.GaugeVec
)

func init() {
	StartGaugeVector()
}

// periodic loop thats wiping the vector.. maybe every hour
func StartGaugeVector() {
	MetricVecs = map[string]prometheus.GaugeVec{}

	go func() {
		for {
			for _, val := range MetricVecs {
				val.Reset()
			}
			time.Sleep(time.Duration(sleepytime) * time.Second)
		}
	}()
}

// Set exports a Prometheus Gauge
func SetMetric(name string, labels map[string]string, value float64) error {
	// maybe create a new gauge vector here to hold the unique name
	gv, ok := MetricVecs[name]
	if ok {
		gv.With(labels).Set(value)
		return nil
	}

	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}

	// the vector does not exist, create it, register and then add this gauge metric to the gauge vector
	gv = *prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: name}, keys)
	customMetrics.Registry.MustRegister(gv)
	MetricVecs[name] = gv

	return nil
}
