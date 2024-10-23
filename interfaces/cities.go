package interfaces

import (
	"context"

	"github.com/shivamk2406/challenge2016/models"
)

type CityService interface {
	GetAllCities(ctx context.Context) ([]*models.City, error)
}