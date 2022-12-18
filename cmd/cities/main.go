package main

import (
	"github.com/admsvist/go-att/internal/pkg/app"
	"log"
)

func main() {
	a := app.New()

	err := a.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
