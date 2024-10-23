package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shivamk2406/challenge2016/models"
	"github.com/shivamk2406/challenge2016/util/csv"
)

type DistributorRepo struct {
	Distributors map[string]*models.Distributor
	CountryMap   map[string]*models.City
	ProvinceMap  map[string]*models.City
	CityMap      map[string]*models.City
}

type API interface {
	AddDistributor(ctx context.Context, distributor *models.Distributor) (*models.Distributor, error)
	GetDistributor(ctx context.Context, name string) (*models.Distributor, error)
	GetLocationByCity(ctx context.Context, city string) *models.City
	GetLocationByProvince(ctx context.Context, province string) *models.City
	GetLocationByCountry(ctx context.Context, country string) *models.City
}

func NewDistributorRepo() API {

	cities, err := GetAllCities(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	countryMap := make(map[string]*models.City)
	provinceMap := make(map[string]*models.City)
	cityMap := make(map[string]*models.City)

	for _, city := range cities {
		countryMap[strings.ToUpper(city.CountryName)] = city
		provinceMap[strings.ToUpper(city.ProvinceName)] = city
		cityMap[strings.ToUpper(city.CityName)] = city
	}

	return &DistributorRepo{
		Distributors: make(map[string]*models.Distributor),
		CountryMap:   countryMap,
		CityMap:      cityMap,
		ProvinceMap:  provinceMap,
	}
}

func (d *DistributorRepo) AddDistributor(ctx context.Context, distributor *models.Distributor) (*models.Distributor, error) {
	d.Distributors[strings.ToUpper(distributor.Name)] = distributor
	return d.GetDistributor(ctx, strings.ToUpper(distributor.Name))
}

func (d *DistributorRepo) GetDistributor(ctx context.Context, name string) (*models.Distributor, error) {

	name = strings.ToUpper(name)

	distributor, ok := d.Distributors[name]
	if !ok {
		return nil, fmt.Errorf("no such distributor exist with name %s", name)
	}

	return distributor, nil
}

func (d *DistributorRepo) GetLocationByCity(ctx context.Context, city string) *models.City {

	name := strings.ToUpper(city)

	loc, ok := d.CityMap[name]
	if !ok {
		return nil
	}
	return loc
}

func (d *DistributorRepo) GetLocationByProvince(ctx context.Context, province string) *models.City {

	name := strings.ToUpper(province)

	loc, ok := d.ProvinceMap[name]
	if !ok {
		return nil
	}
	return loc
}

func (d *DistributorRepo) GetLocationByCountry(ctx context.Context, country string) *models.City {

	name := strings.ToUpper(country)

	loc, ok := d.CountryMap[name]
	if !ok {
		return nil
	}
	return loc
}

func GetAllCities(ctx context.Context) ([]*models.City, error) {

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
