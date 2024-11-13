package dependencies

import "parser/internal/service/emex"

type dependencies struct {
	emexService emex.Service
}



type Dependencies interface {
	FillDeps()
	Run()
}

func NewDependencies() Dependencies {
	return &dependencies{}

}

func (d *dependencies) FillDeps(){
	d.NewEmexService()

}
func (d *dependencies) Run() {
	d.emexService.ParseData()

}
