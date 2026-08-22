package autopiter

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	PauseMs           int    `yaml:"pause_ms"`
	RequestTimeoutSec int    `yaml:"request_timeout_sec"`
	PartsFile         string `yaml:"parts_file"`
	OutputCSV         string `yaml:"output_csv"`
	SkippedCSV        string `yaml:"skipped_csv"`
	NotFoundCSV       string `yaml:"not_found_csv"`
	LogFile           string `yaml:"log_file"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	cfg.applyDefaults()
	return &cfg, nil
}

func (c *Config) applyDefaults() {
	if c.PauseMs <= 0 {
		c.PauseMs = 200
	}
	if c.RequestTimeoutSec <= 0 {
		c.RequestTimeoutSec = 5
	}
	if c.PartsFile == "" {
		c.PartsFile = "internal/files/parts.csv"
	}
	if c.OutputCSV == "" {
		c.OutputCSV = "autopiter.csv"
	}
	if c.SkippedCSV == "" {
		c.SkippedCSV = "autopiter_skipped.csv"
	}
	if c.NotFoundCSV == "" {
		c.NotFoundCSV = "autopiter_notfound.csv"
	}
	if c.LogFile == "" {
		c.LogFile = "autopiter.log"
	}
}
