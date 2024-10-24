package cmd

import (
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
	clisvc:=service.NewCLIService(svc)
	clisvc.RenderCli()
	return nil
}