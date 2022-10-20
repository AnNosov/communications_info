package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HttpServer `yaml:"http_server"`
	HttpClient `yaml:"http_client"`
	FilePath   `yaml:"file_path"`
}

type HttpServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type HttpClient struct {
	MMSHost      string `yaml:"MMS_host"`
	MMSPort      string `yaml:"MMS_port"`
	MethodPath   string `yaml:"MMS_method_path"`
	SupportHost  string `yaml:"support_host"`
	SupportPort  string `yaml:"support_port"`
	SupportPath  string `yaml:"support_method_path"`
	IncidentHost string `yaml:"incident_host"`
	IncidentPort string `yaml:"incident_port"`
	IncidentPath string `yaml:"incident_method_path"`
}

type FilePath struct {
	SmsFilePath          string `yaml:"sms_file_path"`
	SmsFileSeparator     string `yaml:"sms_file_separator"`
	VoiceFilePath        string `yaml:"voice_file_path"`
	VoiceFileSeparator   string `yaml:"voice_file_separator"`
	EmailFilePath        string `yaml:"email_file_path"`
	EmailFileSeparator   string `yaml:"email_file_separator"`
	BillingFilePath      string `yaml:"billing_file_path"`
	BillingFileSeparator string `yaml:"billing_file_separator"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	content, err := os.ReadFile(filepath.Join("config", "config.yaml"))
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal data: %w", err)
	}

	return cfg, nil
}
