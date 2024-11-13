package dependencies

import (
	"parser/internal/service/emex"
	"parser/internal/utils"
)

func (d *dependencies) NewEmexService() emex.Service {
	if d.emexService == nil {
		d.emexService = emex.NewService([]string{"188.119.126.60:9375:a704Dn:gYypM7",
			"194.67.215.106:9917:a704Dn:gYypM7",
			"194.67.216.67:9018:a704Dn:gYypM7",
			"194.67.215.217:9294:a704Dn:gYypM7",
			"194.67.216.46:9591:a704Dn:gYypM7",
			"194.67.213.87:9479:a704Dn:gYypM7",
			"194.67.215.43:9073:a704Dn:gYypM7",
			"194.67.215.90:9150:a704Dn:gYypM7",
			"194.67.215.77:9579:a704Dn:gYypM7",
			"194.67.214.179:9692:a704Dn:gYypM7",
		}, utils.OpenParts("internal/files/parts.csv"), []string{
			"28440", "29629", "32122", "33728", "38920", "38900", "39326", "27269", "23087", "32452",
			"28928", "36383", "31231", "38754", "31026", "11474", "34120", "32492",
		})
	}
	return d.emexService

}
