package client

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CloudResourceSpec struct {}

type CloudResourceService struct {}

var _ CloudResourceOperatorService = (*CloudResourceService)(nil)

func NewCloudResourceService(spec CloudResourceSpec) *CloudResourceService {
	return &CloudResourceService{}
}

type CloudResourceOperatorService interface {
	ReconcileStrategyMaps(ctx context.Context, client client.Client, timeConfig *StrategyTimeConfig, tier, namespace string) error
}




