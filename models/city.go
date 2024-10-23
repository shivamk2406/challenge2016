package models

//City Code,Province Code,Country Code,City Name,Province Name,Country Name

type City struct {
	CityCode     string `json:"City Code,omitempty"`
	ProvinceCode string `json:"Province Code,omitempty"`
	CountryCode  string `json:"Country Code,omitempty"`
	CityName     string `json:"City Name,omitempty"`
	ProvinceName string `json:"Province Name,omitempty"`
	CountryName  string `json:"Country Name,omitempty"`
}
