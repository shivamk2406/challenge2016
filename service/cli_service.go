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
	state     string
	textInput string
	areaInput int
)

const (
	suceesDistributorMessage = "Distributor %s added suceessfully!!!"
	suceesDistributorTitle   = "Success!!!!"
	failureDistributorTitle  = "Failure!!!!"
	successFoundTitle="Found"
	failureFoundTitle="Not Found"
	failureFoundMessage="No such distributor exists with name %s"
	successFoundMessage="Distributor %s found"
	permissionSuccessTitle="Allowed"
	permissionsSuccessMessage="Distributor %s is allowed in given location" 
	permissionsFailureTitle="Not Allowed"
	permissionsFailureMessage="Distributor %s is not allowed in given location"
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
			fmt.Println(err)
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
		fmt.Println(err)
	}
	var excludeAreaResp []models.City

	locationArray := strings.Split(locationString, "-")

	excludeAreaResp = handleArea(locationArray, excludeAreaResp)

	isAllowed := c.distributorService.CheckDistributorPermissions(context.Background(), &models.Permission{
		Name: distributor,
		City: &excludeAreaResp[0],
	})



	if isAllowed{
		showDialogPrompt(permissionSuccessTitle,fmt.Sprintf(permissionsSuccessMessage,distributor))
	}
	return
}

func (c *cliService) triggerDistributorDetailsForm() {
	var distributor *models.Distributor
	var name string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter Distributor Name ").Value(&name),
		))

	err := form.Run()
	if err != nil {
		fmt.Println(err)
	}

	distributor, err = c.distributorService.GetDistributor(context.Background(), name)
	if err != nil {
		fmt.Println(err)
		showDialogPrompt(failureFoundTitle,fmt.Sprintf(failureFoundMessage,name))
	}
	showDialogPrompt(successFoundTitle,fmt.Sprintf(successFoundMessage,distributor.Name))
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
		fmt.Println(err)
	}

	var parent string

	parentInputForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("enter parent for distributor").Value(&parent),
		))

	err = parentInputForm.Run()
	if err != nil {
		fmt.Println(err)
	}

	includeArea, excludeArea := triggerAreaInputForm()

	distributor.IncludedArea = includeArea
	distributor.ExcludedArea = excludeArea
	distributor.Parent = parent

	_, err = c.distributorService.CreateDistributor(context.Background(), &distributor)
	if err != nil {
		fmt.Println(err)
		showDialogPrompt(failureDistributorTitle, fmt.Sprintf(err.Error()))
	}

	showDialogPrompt(suceesDistributorTitle, fmt.Sprintf(suceesDistributorMessage, distributor.Name))

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

func triggerAreaInputForm() ([]models.City, []models.City) {

	var includeAreaResp []models.City
	var excludeAreaResp []models.City
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
	for _, city := range includeAreaResp {
		fmt.Println(city)
	}
	return includeAreaResp, excludeAreaResp

}

func handleArea(includedAreaArray []string, resp []models.City) []models.City {
	switch len(includedAreaArray) {
	case 1:
		resp = append(resp, models.City{
			CountryName: strings.ToUpper(includedAreaArray[0]),
		})

		fmt.Println(includedAreaArray[0], len(includedAreaArray))
		break
	case 2:
		resp = append(resp, models.City{
			CountryName:  strings.ToUpper(includedAreaArray[1]),
			ProvinceName: strings.ToUpper(includedAreaArray[0]),
		})
		fmt.Println(includedAreaArray[0], len(includedAreaArray))
	case 3:
		resp = append(resp, models.City{
			CountryName:  strings.ToUpper(includedAreaArray[2]),
			ProvinceName: strings.ToUpper(includedAreaArray[1]),
			CityName:     strings.ToUpper(includedAreaArray[0]),
		})

		fmt.Println(includedAreaArray[0], len(includedAreaArray))
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

func triggerParentForm() string {
	var resp string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("enter parent for distributor").Value(&resp),
		))
	err := form.Run()
	if err != nil {
		fmt.Println(err)
	}

	return resp
}
