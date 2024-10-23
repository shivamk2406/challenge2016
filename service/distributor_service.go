package service

import (
	"context"
	"errors"
	"strings"

	"github.com/shivamk2406/challenge2016/interfaces"
	"github.com/shivamk2406/challenge2016/models"
	"github.com/shivamk2406/challenge2016/repo"
)

type Service struct {
	repo repo.API
}

func NewDistributorService(repo repo.API) interfaces.DistributorService {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateDistributor(ctx context.Context, distributor *models.Distributor) (*models.Distributor, error) {

	distributorRes, err := s.repo.GetDistributor(ctx, distributor.Name)
	if err != nil {
		return nil, err
	}

	distributor.ExcludedArea=s.populateLocation(ctx, distributor.ExcludedArea)
	distributor.IncludedArea=s.populateLocation(ctx, distributor.IncludedArea)

	if distributor.Parent != nil {
		isAllowed := s.checkParentLocationPermissions(ctx, distributor)
		if !isAllowed {
			return nil, errors.New("parent distributor is not authorized for this location")
		}

		upperCaseName := strings.ToUpper(*distributor.Parent)
		distributor.Parent = &upperCaseName
	}

	distributorRes, err = s.repo.AddDistributor(ctx, distributorRes)
	if err != nil {
		return nil, err
	}
	return distributorRes, nil

}

func (s *Service) GetDistributor(ctx context.Context, name string) (*models.Distributor, error) {
	return s.repo.GetDistributor(ctx, name)
}

func (s *Service) CheckDistributorPermissions() bool {
	return false
}

func getLocationPermission(distributorLoc []models.City, loc models.City) bool {
	// country, province , city
	for _, dLoc := range distributorLoc {
		if strings.EqualFold(dLoc.CountryName, loc.CountryName) {
			if strings.EqualFold(dLoc.ProvinceName, loc.ProvinceName) || strings.EqualFold(dLoc.ProvinceName, "") {
				if strings.EqualFold(dLoc.CityName, loc.CityName) || strings.EqualFold(dLoc.CityName, "") {
					return true
				}
			}
		}
	}
	return false
}

func (s *Service) checkParentLocationPermissions(ctx context.Context, distributor *models.Distributor) bool {

	parent, err := s.repo.GetDistributor(ctx, strings.ToUpper(*&distributor.Name))
	if err == nil {
		return false
	}

	for _, location := range distributor.IncludedArea {
		isExcludedArea := getLocationPermission(parent.ExcludedArea, location)
		isIncludedArea := getLocationPermission(parent.IncludedArea, location)

		if !isIncludedArea || isExcludedArea {
			return false
		}
	}

	return true
}

func(s *Service) populateLocation(ctx context.Context ,loc  []models.City) []models.City{
	for i := range loc{
		if loc[i].CityName != ""{
			loc[i] = *s.repo.GetLocationByCity(ctx ,loc[i].CityName)
		}else if loc[i].ProvinceName != ""{
			loc[i] = *s.repo.GetLocationByProvince(ctx, loc[i].ProvinceName)
		}else if loc[i].CountryName != ""{
			loc[i] = *s.repo.GetLocationByCountry(ctx, loc[i].CountryName)
		}
	}

	return loc
}