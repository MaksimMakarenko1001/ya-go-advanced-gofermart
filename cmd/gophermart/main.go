package main

import (
	"log"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/di"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	di := di.DI{}
	if err := di.Init("GOFERMART_"); err != nil {
		return err
	}

	return di.Start()
}
