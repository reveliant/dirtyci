package server

import (
	"path/filepath"
	"io/ioutil"
	"strings"
	"encoding/json"

	"gopkg.in/yaml.v2"
	"github.com/BurntSushi/toml"
)

type Config struct {
	PluginsDir		string             `yaml:"pluginsDir"`
	Defaults		Repository
	Repositories	Repositories
	Plugins			map[string]string
}

func NewConfig() *Config {
	return new(Config)
}

func (config *Config) SetDefaults(defaults Repository) {
	config.PluginsDir = "plugins"
	config.Defaults = defaults
}

func (config *Config) Load(filename string) error {
	// Select configuration file name and type
	var format = strings.TrimPrefix(filepath.Ext(filename), ".")

	// Read configuration file
	var file, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Read configuration
	switch format {
		case "json":
			err = json.Unmarshal(file, config)
		case "yaml", "yml":
			err = yaml.Unmarshal(file, config)
		case "toml":
			err = toml.Unmarshal(file, config)
	}
	if err != nil {
		return err
	}

	config.Repositories.Sort()
	config.Repositories.SetDefaults(config.Defaults)

	return nil
}

func (config *Config) FindRepo(reponame string) *Repository {
	return config.Repositories.Search(reponame)
}
