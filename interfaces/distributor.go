package interfaces

import (
	"context"

	"github.com/shivamk2406/challenge2016/models"
)

type DistributorService interface {
	//CRUD on Distributor
	CreateDistributor(ctx context.Context, distributor *models.Distributor) (*models.Distributor, error)
	GetDistributor(ctx context.Context, name string) (*models.Distributor, error)

	//Applications
	CheckDistributorPermissions() bool
}
