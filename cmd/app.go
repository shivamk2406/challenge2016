package cmd

import (
	"github.com/shivamk2406/challenge2016/repo"
	"github.com/shivamk2406/challenge2016/service"
)


func Start() error {

	repo := repo.NewDistributorRepo()
	svc := service.NewDistributorService(repo)
	clisvc:=service.NewCLIService(svc)
	clisvc.RenderCli()
	return nil
}