package resources

import (
	"fmt"
	"time"

	errorUtil "github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	customMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	DefaultRedisInfoMetricName = "cro_aws_rds_info"
	DefaultPostgresInfoMetricName = "cro_aws_elasticache_info"
)

// Set exports a Prometheus Gauge
func SetGaugeTimestamp(name string, labels map[string]string, epochTimestamp float64) error {

	// create new gauge
	metric := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        name,
			Help:        labels["description"],
			ConstLabels: labels,
		},
	)

	// un register old metric to allow for new update
	if err := customMetrics.Registry.Unregister(metric); err != false {
		return errorUtil.New(fmt.Sprintf("unable to unregister prometheus metric: %s", metric))
	}
	// register new updated metric
	if err := customMetrics.Registry.Register(metric); err != nil {
		return errorUtil.Wrap(err, fmt.Sprintf("unable register metric %s", metric))
	}

	metric.Set(epochTimestamp)
	return nil
}

func SetGaugeCurrentTime(name string, labels map[string]string) error {
	if err := SetGaugeTimestamp(name, labels, float64(time.Now().UnixNano()) / 1e9); err != nil {
		return errorUtil.Wrap(err,"unable to set gauge")
	}
	return nil
}