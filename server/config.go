package server

import (
	"os"

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

func GetConfigSettings() (*Config, error) {
	config, err := readConfigFile("config-local.yaml")
	if err != nil {
		Log.WithFields(logrus.Fields{
			"Error": "Unable to load config from yaml file",
		}).Fatal(err.Error())
		return nil, err
	}
	return config, nil
}
