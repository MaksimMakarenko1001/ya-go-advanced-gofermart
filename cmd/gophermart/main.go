package main

import (
	"log"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/di"
)

func main() {
	log.Println("app starting")

	if err := run(); err != nil {
		panic(err)
	}

	log.Println("app stoped")
}

func run() error {
	di := di.DI{}
	di.Init("GOFERMART_")

	return di.Start()
}
