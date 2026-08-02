package autopiter

import (
	"strings"

	autopiterconfig "parser/internal/config/autopiter"
	"parser/internal/dependencies"
	"parser/internal/utils"
)

type App interface {
	Start()
}

type app struct {
	cfg  *autopiterconfig.Config
	deps dependencies.AutopiterDependencies
}

func NewApp(cfg *autopiterconfig.Config) App {
	return &app{
		cfg:  cfg,
		deps: dependencies.NewAutopiterDependencies(cfg),
	}
}

func readProxies() []string {
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

func (a *app) Start() {
	proxies := readProxies()
	parts := utils.OpenAutopiterParts(a.cfg.PartsFile)
	a.deps.FillDeps(proxies, parts)
	a.deps.Run()
}
