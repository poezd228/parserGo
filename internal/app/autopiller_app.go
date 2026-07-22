package app

import (
	"fmt"
	"strings"

	"parser/internal/domain"
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
	fmt.Println([]domain.Part{parts[0]})
	a.service = autopiller.NewService(proxies, []domain.Part{parts[0]})
	a.service.ParseData()
}
