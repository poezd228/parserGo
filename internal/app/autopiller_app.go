package app

import (
	"fmt"
	"strings"

	"parser/internal/service/autopiller"
	"parser/internal/utils"
)

type AutopillerApp interface {
	Start()
}

type autopillerApp struct {
	service autopiller.Service
}

func NewAutopillerApp() AutopillerApp {
	return &autopillerApp{}
}

func readAutopillerProxies() []string {
	proxies, err := utils.ReadProxies()
	if err != nil {
		return nil
	}

	var result []string
	for _, proxy := range proxies {
		if trimmed := strings.TrimSpace(proxy); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

func (a *autopillerApp) Start() {
	proxies := readAutopillerProxies()
	parts := utils.OpenAutopiterParts("internal/files/parts.csv")
	fmt.Printf("loaded %d parts from csv, search by partnumber (first=%q)\n", len(parts), parts[0].PartNumber)
	a.service = autopiller.NewService(proxies, parts)
	a.service.ParseData()
}
