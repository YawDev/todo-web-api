package server

import (
	"os"
	"strings"
	l "todo-web-api/loggerutils"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App        App        `yaml:"app"`
	Database   Database   `yaml:"database"`
	Swagger    Swagger    `yaml:"swagger"`
	APIConfig  APIConfig  `yaml:"api"`
	CORSConfig CORSConfig `yaml:"cors"`
}

type App struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Port        string `yaml:"port"`
	Host        string `yaml:"host"`
	BindAddress string `yaml:"bind_address"`
}

type Database struct {
	Driver    string `yaml:"driver"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Name      string `yaml:"name"`
	SSLMode   string `yaml:"ssl_mode"`
	UseSQLite bool   `yaml:"useSQLite"`
}

type Swagger struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
	DocPath string `yaml:"doc_path"`
}

type APIConfig struct {
	BaseURL string `yaml:"base_url"`
	Timeout int    `yaml:"timeout"`
}

type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
}

func readConfigFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetConfigSettings() *Config {
	// CONFIG_FILE selects which config to load (e.g. config.production.yaml in
	// deployment); defaults to config.yaml for local dev.
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.yaml"
	}
	config, err := readConfigFile(configFile)
	if err != nil {
		l.Log.WithFields(logrus.Fields{
			"Error": "Unable to load config from yaml file",
			"File":  configFile,
		}).Fatal(err.Error())
		return nil
	}
	applyEnvOverrides(config)
	return config
}

// applyEnvOverrides lets the deployment environment override config.yaml
// values without rebuilding the image. Anything unset falls back to the file.
func applyEnvOverrides(config *Config) {
	// APP_ENV marks the deployment (e.g. "production"); anything other than
	// "local-development" disables dev-only behavior like auto-opening Swagger.
	if env := os.Getenv("APP_ENV"); env != "" {
		config.App.Environment = env
	}
	// Fly (and most platforms) inject the listening port via PORT.
	if port := os.Getenv("PORT"); port != "" {
		config.App.Port = port
	}
	// Comma-separated list of allowed CORS origins, e.g.
	// "https://todo-manager-yaw-dev.vercel.app".
	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		parts := strings.Split(origins, ",")
		trimmed := make([]string, 0, len(parts))
		for _, p := range parts {
			if v := strings.TrimSpace(p); v != "" {
				trimmed = append(trimmed, v)
			}
		}
		config.CORSConfig.AllowedOrigins = trimmed
	}
}
