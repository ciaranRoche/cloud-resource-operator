package resources

import (
	"fmt"
	"time"

	"github.com/integr8ly/cloud-resource-operator/version"
	prometheus "github.com/prometheus/client_golang/prometheus"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	customMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
	metricsstore "k8s.io/kube-state-metrics/pkg/metrics_store"
)

var (
	log        = logf.Log.WithName("custom_metrics")
	sleepytime = 5

	operatorVersion = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        "cro_version_info",
			Help:        "CRO operator information",
			ConstLabels: prometheus.Labels{"operator_version": version.Version},
		},
	)

	// INTLY-4669
	upcommingRDSMaintenance = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        "cro_aws_rds_scheduled_maintenance_window",
			Help:        "AWS RDS Resources with scheduled maintenance windows managed by the CRO operator",
			ConstLabels: prometheus.Labels{"aws_region": "eu-west-1"},
		},
	)

	upcommingElastiCacheMaintenance = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        "cro_aws_elasticache_scheduled_maintenance_window",
			Help:        "AWS ElastiCache Resources with scheduled maintenance windows managed by the CRO operator",
			ConstLabels: prometheus.Labels{"aws_region": "eu-west-1"},
		},
	)
)

func addCustomMetrics() {
	customMetrics.Registry.MustRegister(operatorVersion)
	customMetrics.Registry.MustRegister(upcommingRDSMaintenance)
	customMetrics.Registry.MustRegister(upcommingElastiCacheMaintenance)
}

func reconcileMetrics() {
	log.Info("custom_metrics reconcileMetrics")

	// datastructure array/hash/whatever of metrics
	// calls functions to update this list of metrics
	// using MetricsStore

	// Reset() 
}

func StartMetricsLoop() {
	go func() {
		for {
			// put the update logic in here
			log.Info(fmt.Sprintf("Go to sleep for: %d seconds", sleepytime))
			reconcileMetrics()
			time.Sleep(time.Duration(sleepytime) * time.Second)
		}
	}()
}
