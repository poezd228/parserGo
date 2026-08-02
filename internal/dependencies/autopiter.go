package dependencies

import (
	autopiterconfig "parser/internal/config/autopiter"
	"parser/internal/domain"
	"parser/internal/service/autopiter"
)

type AutopiterDependencies interface {
	FillDeps(proxies []string, parts []domain.Part)
	Run()
}

type autopiterDependencies struct {
	cfg             *autopiterconfig.Config
	autopiterService autopiter.Service
}

func NewAutopiterDependencies(cfg *autopiterconfig.Config) AutopiterDependencies {
	return &autopiterDependencies{cfg: cfg}
}

func (d *autopiterDependencies) FillDeps(proxies []string, parts []domain.Part) {
	d.NewAutopiterService(proxies, parts)
}

func (d *autopiterDependencies) NewAutopiterService(proxies []string, parts []domain.Part) autopiter.Service {
	if d.autopiterService == nil {
		d.autopiterService = autopiter.NewService(d.cfg, proxies, parts)
	}
	return d.autopiterService
}

func (d *autopiterDependencies) Run() {
	d.autopiterService.ParseData()
}
