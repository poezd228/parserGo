package dependencies

import (
	"parser/internal/domain"
	"parser/internal/service/emex"
)

type dependencies struct {
	emexService emex.Service
}

type Dependencies interface {
	FillDeps(proxies []string, parts []domain.Part, locations []string)
	Run()
}

func NewDependencies() Dependencies {
	return &dependencies{}

}

func (d *dependencies) FillDeps(proxies []string, parts []domain.Part, locations []string) {
	d.NewEmexService(proxies, parts, locations)

}
func (d *dependencies) Run() {
	d.emexService.ParseData()

}
