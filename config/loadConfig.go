package config

import (
	_ "embed"

	"github.com/Anacardo89/lenic_api/internal/server"
	"github.com/Anacardo89/lenic_api/pkg/db"
	"gopkg.in/yaml.v3"
)

//go:embed dbConfig.yaml
var dbYaml []byte

//go:embed serverConfig.yaml
var serverYaml []byte

func LoadDBConfig() (*db.Config, error) {
	var config db.Config
	err := yaml.Unmarshal(dbYaml, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadServerConfig() (*server.Config, error) {
	var config server.Config
	err := yaml.Unmarshal(serverYaml, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
