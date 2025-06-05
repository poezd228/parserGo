package app

import (
	"parser/internal/dependencies"
	"parser/internal/utils"
)

type app struct {
	deps dependencies.Dependencies
}
type App interface {
	Start()
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
	// locations := []string{
	// 	"28440", "29629", "32122", "33728", "38920", "38900", "39326", "27269", "23087", "32452",
	// 	"28928", "36383", "31231", "38754", "31026", "11474", "34120", "32492",
	// }
	locations := []string{
		"16733", "29435", "30254", "638",
	}
	locationsCoords := make(map[string][]string)
	locationsCoords["16733"] = []string{"131.912", "43.121"}
	locationsCoords["30254"] = []string{"131.949", "43.1474"}

	a.deps.FillDeps(proxies, parts, locations, locationsCoords)
	a.deps.Run()

}
