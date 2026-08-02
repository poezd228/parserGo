package emex

import (
	"parser/internal/dependencies"
	"parser/internal/utils"
)

type App interface {
	Start()
}

type app struct {
	deps dependencies.Dependencies
}

func NewApp() App {
	return &app{deps: dependencies.NewDependencies()}
}

func (a *app) Start() {
	proxies, err := utils.ReadProxies()
	if err != nil {
		panic(err)
	}
	parts := utils.OpenParts("internal/files/parts.csv")
	locations := []string{
		"16733", "29435", "30254", "638",
	}
	locationsCoords := make(map[string][]string)
	locationsCoords["16733"] = []string{"131.912", "43.121"}
	locationsCoords["30254"] = []string{"131.949", "43.1474"}

	a.deps.FillDeps(proxies, parts, locations, locationsCoords)
	a.deps.Run()
}
