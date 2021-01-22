package ircd

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// Config loaded from config.yaml
type Config struct {
	Name string
	Motd string

	Address string

	Channels []string
}

// LoadConfiguration uses a path to load and parse config.yaml
func LoadConfiguration(path string) (*Config, error) {
	config := Config{}

	// read config file
	data, e := ioutil.ReadFile(path)
	if e != nil {
		return nil, e
	}

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// config.required("Name") TODO required values
	// config.required("Address")
	// config.required("Channels")

	return &config, nil
}

// StoreConfiguration saves the current config to config.yaml
func (config Config) StoreConfiguration() error {

	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("config.yaml", data, 0644)
	if err != nil {
		return err
	}

	return nil
}
