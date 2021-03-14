package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represents the top-level structure of the configuration file.
type Config struct {
	Scanner *ScannerConfig `yaml:"scanner"`
	HTTP    *HTTPConfig    `yaml:"http"`
	WebDAV  *WebDAVConfig  `yaml:"webdav"`
}

// ScannerConfig represents the configuration for the scanner, i.e. the device that's
// scanning documents.
type ScannerConfig struct {
	DeviceName string `yaml:"device_name"`
	Mode       string `yaml:"mode"`
	PreviewRes int    `yaml:"preview_res"`
	ScanRes    int    `yaml:"scan_res"`
}

// HTTPConfig represents the configuration for the HTTP server used to preview scans and
// control the scan of documents by the scanner.
type HTTPConfig struct {
	Address string `yaml:"bind_addr"`
	Port    string `yaml:"bind_port"`
}

// WebDAVConfig represents the configuration required to connect to the WebDAV server and
// upload scanned documents there.
type WebDAVConfig struct {
	RootURL    string `yaml:"root_url"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	UploadPath string `yaml:"upload_path"`
}

// NewConfig parses the configuration file at the given path.
func NewConfig(path string) (*Config, error) {
	configWithDefaults := &Config{
		Scanner: &ScannerConfig{},
		HTTP: &HTTPConfig{
			Address: "127.0.0.1",
			Port:    "8080",
		},
		WebDAV: &WebDAVConfig{},
	}

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(raw, configWithDefaults); err != nil {
		return nil, err
	}

	return configWithDefaults, nil
}
