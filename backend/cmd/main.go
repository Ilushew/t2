package main

import (
	"log"

	"github.com/ilushew/udmurtia-trip/backend/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	a.Run()
}
