package service

import (
	"context"
	"encoding/json"

	"github.com/shivamk2406/challenge2016/interfaces"
	"github.com/shivamk2406/challenge2016/models"
	"github.com/shivamk2406/challenge2016/util/csv"
)

type cityService struct {
}

func NewCityService(fileName string) interfaces.CityService {
	return &cityService{}
}

func (c *cityService) GetAllCities(ctx context.Context) ([]*models.City, error) {

	rowentries, err := csv.ReadCsvFile(ctx, "cities.csv")
	if err != nil {
		return nil, err
	}

	cities := make([]*models.City, 0)

	bytes, err := json.Marshal(rowentries)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &cities)
	if err != nil {
		return nil, err
	}

	return cities, nil
}
