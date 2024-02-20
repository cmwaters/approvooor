package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const Name = "config.toml"

type Config struct {
	Address string `toml:"address"`
}

func Default() *Config {
	return &Config{
		Address: ":8080",
	}
}

func Load(dir string) (*Config, error) {
	var config Config
	path := filepath.Join(dir, Name)
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func Exists(dir string) bool {
	path := filepath.Join(dir, Name)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (c *Config) Save(dir string) error {
	path := filepath.Join(dir, Name)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(c); err != nil {
		return err
	}
	return nil
}
