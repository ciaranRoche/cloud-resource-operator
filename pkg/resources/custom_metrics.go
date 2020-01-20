package resources

import (
	"github.com/integr8ly/cloud-resource-operator/version"
	prometheus "github.com/prometheus/client_golang/prometheus"
	customMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
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

func AddCustomMetrics() {
	customMetrics.Registry.MustRegister(operatorVersion)
	customMetrics.Registry.MustRegister(upcommingRDSMaintenance)
	customMetrics.Registry.MustRegister(upcommingElastiCacheMaintenance)
}
