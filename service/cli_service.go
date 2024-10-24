package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/shivamk2406/challenge2016/interfaces"
	"github.com/shivamk2406/challenge2016/models"
)

var (
	input     int
	areaInput int
)

const (
	suceesDistributorMessage  = "Distributor %s added suceessfully!!!"
	suceesDistributorTitle    = "Success!!!!"
	failureDistributorTitle   = "Failure!!!!"
	successFoundTitle         = "Found"
	failureFoundTitle         = "Not Found"
	failureFoundMessage       = "No such distributor exists with name %s"
	successFoundMessage       = "Distributor %s found"
	permissionSuccessTitle    = "Allowed"
	permissionsSuccessMessage = "Distributor %s is allowed in given location"
	permissionsFailureTitle   = "Not Allowed"
	permissionsFailureMessage = "Distributor %s is not allowed in given location"
	genericErrorTitle         = "Error"
)

type cliService struct {
	distributorService interfaces.DistributorService
}

func NewCLIService(svc interfaces.DistributorService) interfaces.CLIRender {
	return &cliService{
		distributorService: svc,
	}
}

func (c *cliService) RenderCli() {

	for {
		form := getStarteForm()
		err := form.Run()
		if err != nil {
			showDialogPrompt(genericErrorTitle, err.Error())
		}

		switch input {
		case 1:
			c.triggerDistributorInputForm()
		case 2:
			c.triggerDistributorDetailsForm()
		case 3:
			c.triggerDistributorPermissionsForm()
		case 4:
			break
		}

		if input == 4 {
			break
		}
	}

}

func getStarteForm() *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Options(
					huh.NewOption("1. Add Distributor", 1),
					huh.NewOption("2. Get Distributor", 2),
					huh.NewOption("3. Get Distributor Permission", 3),
					huh.NewOption("4. Exit", 4),
				).
				Value(&input).
				Title("Select operation").
				Height(15),
		))
}

func (c *cliService) triggerDistributorPermissionsForm() {
	var distributor string
	var locationString string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter Distributor Name ").Value(&distributor),
			huh.NewInput().Title("Enter Location Name ").Value(&locationString),
		))

	err := form.Run()
	if err != nil {
		showDialogPrompt(genericErrorTitle, err.Error())
	}
	var excludeAreaResp []models.City

	locationArray := strings.Split(locationString, "-")

	excludeAreaResp = handleArea(locationArray, excludeAreaResp)

	isAllowed := c.distributorService.CheckDistributorPermissions(context.Background(), &models.Permission{
		Name: distributor,
		City: &excludeAreaResp[0],
	})

	if isAllowed {
		showDialogPrompt(permissionSuccessTitle, fmt.Sprintf(permissionsSuccessMessage, distributor))
	}
	showDialogPrompt(permissionsFailureTitle, fmt.Sprintf(permissionsFailureMessage, distributor))
	return
}

func (c *cliService) triggerDistributorDetailsForm() {
	var name string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter Distributor Name ").Value(&name),
		))

	err := form.Run()
	if err != nil {
		showDialogPrompt(genericErrorTitle, err.Error())
	}

	_, err = c.distributorService.GetDistributor(context.Background(), name)
	if err != nil {
		fmt.Println(err)
		showDialogPrompt(failureFoundTitle, fmt.Sprintf(failureFoundMessage, name))
	}
	showDialogPrompt(successFoundTitle, fmt.Sprintf(successFoundMessage, name))
	return
}

func (c *cliService) triggerDistributorInputForm() {
	var distributor models.Distributor

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter Distributor Name ").Value(&distributor.Name),
		))

	err := form.Run()
	if err != nil {
		showDialogPrompt(genericErrorTitle, err.Error())
	}

	var parent string

	triggerParentInputForm(parent, distributor.Name)

	includeArea, excludeArea := triggerAreaInputForm(distributor.Name)

	distributor.IncludedArea = includeArea
	distributor.ExcludedArea = excludeArea
	distributor.Parent = parent


	_, err = c.distributorService.CreateDistributor(context.Background(), &distributor)
	if err != nil {
		showDialogPrompt(failureDistributorTitle, fmt.Sprintf(err.Error()))
	}

	showDialogPrompt(suceesDistributorTitle, fmt.Sprintf(suceesDistributorMessage, distributor.Name))

}

func triggerParentInputForm(parent, distributor string) {
	parentInputForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title(fmt.Sprintf("Enter parent name for distributor %s", distributor)).Value(&parent),
		))

	err := parentInputForm.Run()
	if err != nil {
		showDialogPrompt(genericErrorTitle, err.Error())
	}
}

func showDialogPrompt(title, message string) {
	dialogue := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(title).Description(message),
		))
	err := dialogue.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func triggerAreaInputForm(distributorName string) ([]models.City, []models.City) {

	var includeAreaResp []models.City
	var excludeAreaResp []models.City
	var includeArea string
	var excludeArea string

	for {
		form := getAreaInputForm()
		err := form.Run()
		if err != nil {
			showDialogPrompt(genericErrorTitle, err.Error())
		}
		switch areaInput {
		case 1:
			includeArea = getAreaInputResponseFromUser(fmt.Sprintf("EXCLUDE for %s: format city-province-state",distributorName))
			includedAreaArray := strings.Split(includeArea, "-")
			includeAreaResp = handleArea(includedAreaArray, includeAreaResp)
		case 2:
			excludeArea = getAreaInputResponseFromUser(fmt.Sprintf("INCLUDE for %s: format city-province-state",distributorName))
			excludedAreaArray := strings.Split(excludeArea, "-")
			excludeAreaResp = handleArea(excludedAreaArray, excludeAreaResp)
		case 3:
			break
		}

		if areaInput == 3 {
			break
		}
	}

	return includeAreaResp, excludeAreaResp

}

func handleArea(includedAreaArray []string, resp []models.City) []models.City {
	switch len(includedAreaArray) {
	case 1:
		resp = append(resp, models.City{
			CountryName: strings.ToUpper(includedAreaArray[0]),
		})
		break
	case 2:
		resp = append(resp, models.City{
			CountryName:  strings.ToUpper(includedAreaArray[1]),
			ProvinceName: strings.ToUpper(includedAreaArray[0]),
		})
	case 3:
		resp = append(resp, models.City{
			CountryName:  strings.ToUpper(includedAreaArray[2]),
			ProvinceName: strings.ToUpper(includedAreaArray[1]),
			CityName:     strings.ToUpper(includedAreaArray[0]),
		})
	}
	return resp
}

func getAreaInputForm() *huh.Form {

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Options(huh.NewOption("1. Insert included area of the form of city-province-country", 1),
					huh.NewOption("2. Insert excluded area of the form of city-province-country", 2),
					huh.NewOption("3. Exit", 3)).
				Value(&areaInput).
				Title("Select the operation").
				Height(15).WithTheme(huh.ThemeCatppuccin()),
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
		showDialogPrompt(genericErrorTitle, err.Error())
	}

	return resp
}
