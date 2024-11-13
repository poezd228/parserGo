package app

import "parser/internal/dependencies"

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
	a.deps.FillDeps()
	a.deps.Run()

}
