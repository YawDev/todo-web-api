package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App        App        `yaml:"app"`
	Database   Database   `yaml:"database"`
	Swagger    Swagger    `yaml:"swagger"`
	APIConfig  APIConfig  `yaml:"apiConfig"`
	CORSConfig CORSConfig `yaml:"corsConfig"`
}

type App struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Port        int    `yaml:"port"`
	Host        string `yaml:"host"`
}

type Database struct {
	Driver    string `yaml:"driver"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Name      string `yaml:"name"`
	SSLMode   string `yaml:"ssl_mode"`
	UseSQLite bool   `yaml:"use_sqlite"`
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

func loadConfig(filename string) (*Config, error) {
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
