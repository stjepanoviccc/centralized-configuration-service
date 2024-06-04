package repositories

import (
	"ars_projekat/model"
	"errors"
	"fmt"
)

type ConfigGroupInMemoryRepository struct {
	configGroups map[string]model.ConfigurationGroup
}

func (c ConfigGroupInMemoryRepository) Add(configGroup *model.ConfigurationGroup) error {
	if configGroup == nil {
		return errors.New("cannot add nil configuration group")
	}
	key := fmt.Sprintf("%s/%#v", configGroup.Name, configGroup.Version)
	if _, exists := c.configGroups[key]; exists {
		return errors.New("configuration group already exists")
	}

	c.configGroups[key] = *configGroup
	return nil
}

func (c ConfigGroupInMemoryRepository) Get(name string, version model.Version) (model.ConfigurationGroup, error) {
	key := fmt.Sprintf("%s/%#v", name, version)
	configGroup, ok := c.configGroups[key]
	if !ok {
		return model.ConfigurationGroup{}, errors.New("config not found")
	}
	return configGroup, nil
}

func (c ConfigGroupInMemoryRepository) Delete(configGroup model.ConfigurationGroup) error {
	key := fmt.Sprintf("%s/%#v", configGroup.Name, configGroup.Version)
	_, ok := c.configGroups[key]
	if !ok {
		return errors.New("config not found")
	}
	delete(c.configGroups, key)
	return nil
}

func (c ConfigGroupInMemoryRepository) Save(configGroup *model.ConfigurationGroup) error {
	if configGroup == nil {
		return errors.New("cannot add nil configuration group")
	}
	key := fmt.Sprintf("%s/%#v", configGroup.Name, configGroup.Version)
	c.configGroups[key] = *configGroup

	return nil
}

func NewConfigGroupInMemoryRepository() model.ConfigurationGroupRepository {
	return ConfigGroupInMemoryRepository{
		configGroups: make(map[string]model.ConfigurationGroup),
	}
}
