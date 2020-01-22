package resources

import (
	"fmt"
	"time"

	errorUtil "github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	DefaultRedisInfoMetricName = "cro_aws_elasticache_info"
	DefaultPostgresInfoMetricName = "cro_aws_rds_info"
)

// Set exports a Prometheus Gauge
func SetGaugeTimestamp(name string, labels map[string]string, epochTimestamp float64) error {
	// create new gauge
	//metric := prometheus.NewGauge(
	//	prometheus.GaugeOpts{
	//		Name:        name,
	//		Help:        labels["description"],
	//		ConstLabels: labels,
	//	},
	//)

	//// un register old metric to allow for new update
	//customMetrics.Registry.Unregister(metric)
	//// register new updated metric
	//if err := customMetrics.Registry.Register(metric); err != nil {
	//	return errorUtil.Wrap(err, fmt.Sprintf("unable register metric %s", metric))
	//}


	gaug, err := GetPostgresInfoGauge().GetMetricWithLabelValues("clusterID")
	if err != nil {
		return err
	}
	fmt.Println("found", gaug.Desc())
	gaug.SetToCurrentTime()

	return nil
}

// returns postgres info gauge vector
func GetPostgresInfoGauge() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "worker",
		Subsystem: "jobs",
		Name: DefaultPostgresInfoMetricName,
		Help: "Returning information on active postgres resources",
	},
	[]string{
		"clusterID",
		//"resourceID",
		//"namespace",
		//"instanceID",
		//"status",
	})
}

func SetGaugeCurrentTime(name string, labels map[string]string) error {
	if err := SetGaugeTimestamp(name, labels, float64(time.Now().UnixNano()) / 1e9); err != nil {
		return errorUtil.Wrap(err,"unable to set gauge")
	}
	return nil
}
