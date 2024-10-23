package main

import (
	"log"

	"github.com/shivamk2406/challenge2016/cmd"
)

func main() {
	err := cmd.Start()
	if err != nil {
		log.Println(err)
	}
}