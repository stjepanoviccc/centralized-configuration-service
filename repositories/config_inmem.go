package repositories

import (
	"ars_projekat/model"
	"errors"
	"fmt"
)

type ConfigInMemoryRepository struct {
	configs map[string]model.Configuration
}

func (c ConfigInMemoryRepository) Add(config *model.Configuration) error {
	if config == nil {
		return errors.New("cannot add nil config")
	}

	key := fmt.Sprintf("%s/%#v", config.Name, config.Version)
	if _, exists := c.configs[key]; exists {
		return errors.New("config already exists")
	}
	c.configs[key] = *config
	return nil
}

func (c ConfigInMemoryRepository) Get(name string, version model.Version) (model.Configuration, error) {
	key := fmt.Sprintf("%s/%#v", name, version)
	config, ok := c.configs[key]
	if !ok {
		return model.Configuration{}, errors.New("config not found")
	}
	return config, nil
}

func (c ConfigInMemoryRepository) Delete(config model.Configuration) error {
	key := fmt.Sprintf("%s/%#v", config.Name, config.Version)
	_, ok := c.configs[key]
	if !ok {
		return errors.New("config not found")
	}
	delete(c.configs, key)
	return nil
}

func NewConfigInMemoryRepository() model.ConfigurationRepository {
	return ConfigInMemoryRepository{
		configs: make(map[string]model.Configuration),
	}
}
