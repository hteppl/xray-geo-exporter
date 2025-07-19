package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"time"
)

var (
	LogFile                string
	WebhookURL             string
	Hostname               string
	IpRegex                = regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)`)
	IPTTLMinutes           int
	RateLimitPerMinute     int
	QueueSize              int
	HTTPTimeout            time.Duration
	WebhookTimeout         time.Duration
	CleanupIntervalMinutes int
)

type Config struct {
	LogFile                string `yaml:"LogFile"`
	WebhookURL             string `yaml:"WebhookURL"`
	Hostname               string `yaml:"Hostname"`
	IPTTLMinutes           int    `yaml:"IPTTLMinutes"`
	RateLimitPerMinute     int    `yaml:"RateLimitPerMinute"`
	QueueSize              int    `yaml:"QueueSize"`
	HTTPTimeoutSeconds     int    `yaml:"HTTPTimeoutSeconds"`
	WebhookTimeoutSeconds  int    `yaml:"WebhookTimeoutSeconds"`
	CleanupIntervalMinutes int    `yaml:"CleanupIntervalMinutes"`
}

func LoadConfig(configPath string) error {
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var cfg Config
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return err
	}

	LogFile = cfg.LogFile
	WebhookURL = cfg.WebhookURL

	// Set defaults if not specified
	Hostname = cfg.Hostname
	if Hostname == "" {
		Hostname, err = os.Hostname()
	}

	IPTTLMinutes = cfg.IPTTLMinutes
	if IPTTLMinutes == 0 {
		IPTTLMinutes = 30
	}

	RateLimitPerMinute = cfg.RateLimitPerMinute
	if RateLimitPerMinute == 0 {
		RateLimitPerMinute = 45
	}

	QueueSize = cfg.QueueSize
	if QueueSize == 0 {
		QueueSize = 3000
	}

	httpTimeout := cfg.HTTPTimeoutSeconds
	if httpTimeout == 0 {
		httpTimeout = 5
	}
	HTTPTimeout = time.Duration(httpTimeout) * time.Second

	webhookTimeout := cfg.WebhookTimeoutSeconds
	if webhookTimeout == 0 {
		webhookTimeout = 5
	}
	WebhookTimeout = time.Duration(webhookTimeout) * time.Second

	CleanupIntervalMinutes = cfg.CleanupIntervalMinutes
	if CleanupIntervalMinutes == 0 {
		CleanupIntervalMinutes = 1
	}

	return err
}
