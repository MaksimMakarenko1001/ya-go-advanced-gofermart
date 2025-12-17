package main

import (
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/di"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	di := di.DI{}
	di.Init("GOFERMART_")

	return di.Start()
}
