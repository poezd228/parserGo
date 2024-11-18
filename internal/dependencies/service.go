package dependencies

import (
	"parser/internal/domain"
	"parser/internal/service/emex"
)

func (d *dependencies) NewEmexService(proxies []string, parts []domain.Part, locations []string) emex.Service {
	if d.emexService == nil {
		d.emexService = emex.NewService(proxies, parts, locations)
	}
	return d.emexService

}
