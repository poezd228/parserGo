package main

import (
	"flag"
	"parser/internal/app"
)

func main() {

	var location string

	flag.StringVar(&location, "location", "", "")
	flag.Parse()

	a := app.NewApp()
	a.Start(location)

}
