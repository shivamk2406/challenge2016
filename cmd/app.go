package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/shivamk2406/challenge2016/models"
	"github.com/shivamk2406/challenge2016/repo"
	"github.com/shivamk2406/challenge2016/service"
)

var (
	burger       string
	toppings     []string
	sauceLevel   int
	name         string
	instructions string
	discount     bool
)

var (
	input     int
	state     string
	textInput string
	areaInput int
)

func Start() error {
	repo := repo.NewDistributorRepo()
	svc := service.NewDistributorService(repo)

	for {
		form := getStarteForm()
		err := form.Run()
		if err != nil {
			fmt.Println(err)
		}

		switch input {
		case 1:
			triggerDistributorInputForm()
		case 2:
			triggerDistributorDetailsForm()
		case 3:
			triggerDistributorPermissionsForm()
		case 4:
			break
		}

		if input == 4 {
			break
		}
	}
	includedCity := make([]models.City, 0)
	includedCity = append(includedCity, models.City{
		CountryName: "india",
	})

	excludedCity := make([]models.City, 0)
	excludedCity = append(excludedCity, models.City{
		CityName: "ranchi",
	})

	svc.CreateDistributor(context.Background(), &models.Distributor{
		Name:         "DISTRIBUTOR1",
		IncludedArea: includedCity,
		ExcludedArea: excludedCity,
	})

	fmt.Println(svc.GetDistributor(context.Background(), "DISTRIBUTOR1"))

	fmt.Println(svc.CheckDistributorPermissions(context.Background(), &models.Permission{
		Name: "DISTRIBUTOR1",
		City: &models.City{
			CityName:     "ranchi",
			CountryName:  "india",
			ProvinceName: "jharkhand",
		},
	}))
	return nil
}

func getStarteForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Options(
					huh.NewOption("1. Add Distributor:", 1),
					huh.NewOption("2. Get Distributor:", 2),
					huh.NewOption("3. Get Distributor Permission:", 3),
					huh.NewOption("4. Exit", 4),
				).
				Value(&input).
				Title("Select the operation").
				Height(15),
		))
}

func triggerDistributorPermissionsForm() {
	var distributor models.Distributor
	var locationString string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter Distributor Name ").Value(&distributor.Name),
			huh.NewInput().Title("Enter Location Name ").Value(&locationString),
		))

	err := form.Run()
	if err != nil {
		fmt.Println(err)
	}

	return
}

func triggerDistributorDetailsForm() {
	var distributor models.Distributor
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter Distributor Name ").Value(&distributor.Name),
		))

	err := form.Run()
	if err != nil {
		fmt.Println(err)
	}

	return
}

func triggerDistributorInputForm() {
	var distributor models.Distributor

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter Distributor Name ").Value(&distributor.Name),
		))

	err := form.Run()
	if err != nil {
		fmt.Println(err)
	}

	triggerAreaInputForm()

	return

}

func triggerAreaInputForm() []*models.City {

	var includeAreaResp []*models.City
	var excludeAreaResp []*models.City
	var includeArea string
	var excludeArea string

	for {
		form := getAreaInputForm()
		err := form.Run()
		if err != nil {
			fmt.Println(err)
		}
		switch areaInput {
		case 1:
			includeArea = getAreaInputResponseFromUser("INCLUDE")
			includedAreaArray := strings.Split(includeArea, "-")
			includeAreaResp = handleArea(includedAreaArray, includeAreaResp)
		case 2:
			excludeArea = getAreaInputResponseFromUser("EXCLUDE")
			excludedAreaArray := strings.Split(excludeArea, "-")
			excludeAreaResp = handleArea(excludedAreaArray, excludeAreaResp)
		case 3:
			break
		}

		if areaInput == 3 {
			break
		}
	}

	fmt.Println(includeArea)
	fmt.Println(excludeArea)
	for _, city := range resp {
		fmt.Println(city)  
	}
	return resp

}

func handleArea(includedAreaArray []string, resp []*models.City) []*models.City {
	switch len(includedAreaArray) {
	case 1:
		resp = append(resp, &models.City{
			CountryName: strings.ToUpper(includedAreaArray[0]),
		})

		fmt.Println(includedAreaArray[0],len(includedAreaArray))
		break
	case 2:
		resp = append(resp, &models.City{
			CountryName:  strings.ToUpper(includedAreaArray[1]),
			ProvinceName: strings.ToUpper(includedAreaArray[0]),
		})
		fmt.Println(includedAreaArray[0],len(includedAreaArray))
	case 3:
		resp = append(resp, &models.City{
			CountryName:  strings.ToUpper(includedAreaArray[2]),
			ProvinceName: strings.ToUpper(includedAreaArray[1]),
			CityName:     strings.ToUpper(includedAreaArray[0]),
		})

		fmt.Println(includedAreaArray[0],len(includedAreaArray))
	}
	return resp
}

func getAreaInputForm() *huh.Form {

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Options(huh.NewOption("1. Insert include area", 1),
					huh.NewOption("2. Insert Exclude area", 2),
					huh.NewOption("3. Exit", 3)).
				Value(&areaInput).
				Title("Select the operation").
				Height(15),
		))
}

func getAreaInputResponseFromUser(areaType string) string {
	var resp string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title(fmt.Sprintf("%s", areaType)).Value(&resp),
		))

	err := form.Run()
	if err != nil {
		fmt.Println(err)
	}

	return resp
}
