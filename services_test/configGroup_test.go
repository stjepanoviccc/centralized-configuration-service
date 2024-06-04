package services_test

import (
	"ars_projekat/model"
	"ars_projekat/repositories"
	"ars_projekat/services"
	"context"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConfigurationGroupService_Add(t *testing.T) {
	configGroup := model.ConfigurationGroup{
		Name:    "testGroup",
		Version: model.Version{Major: 1, Minor: 0, Patch: 0},
		Configurations: []model.Configuration{
			{Name: "config1", Version: model.Version{Major: 1, Minor: 0, Patch: 0}, Labels: map[string]string{"label1": "value1"}},
		},
	}

	mockRepo := new(repositories.MockConfigRepository)

	// Set up mock expectations for each configuration in the group
	for _, config := range configGroup.Configurations {
		mockRepo.On("AddGroup", configGroup.Name, model.ToString(configGroup.Version), model.SortLabels(config.Labels), config, mock.Anything).Return(nil)
	}

	service := services.NewConfigurationGroupService(mockRepo, NewTestTracer())

	err := service.Add(configGroup, context.Background())
	assert.NoError(t, err)

	// Check that AddGroup was called with the expected parameters for each configuration
	for _, config := range configGroup.Configurations {
		mockRepo.AssertCalled(t, "AddGroup", configGroup.Name, model.ToString(configGroup.Version), model.SortLabels(config.Labels), config, mock.Anything)
	}
	mockRepo.AssertExpectations(t)
}

func TestConfigurationGroupService_Save(t *testing.T) {
	mockRepo := new(repositories.MockConfigRepository)
	service := services.NewConfigurationGroupService(mockRepo, NewTestTracer())

	configGroup := &model.ConfigurationGroup{
		Name:    "testGroup",
		Version: model.Version{Major: 1, Minor: 0, Patch: 0},
		Configurations: []model.Configuration{
			{Name: "config1", Version: model.Version{Major: 1, Minor: 0, Patch: 0}, Labels: map[string]string{"label1": "value1"}},
		},
	}

	for _, config := range configGroup.Configurations {
		mockRepo.On("AddGroup", configGroup.Name, model.ToString(configGroup.Version), model.SortLabels(config.Labels), config, mock.Anything).Return(nil)
	}

	err := service.Save(configGroup, context.Background())
	assert.NoError(t, err)

	for _, config := range configGroup.Configurations {
		mockRepo.AssertCalled(t, "AddGroup", configGroup.Name, model.ToString(configGroup.Version), model.SortLabels(config.Labels), config, mock.Anything)
	}
	mockRepo.AssertExpectations(t)
}

func TestConfigurationGroupService_Get(t *testing.T) {
	mockRepo := new(repositories.MockConfigRepository)
	service := services.NewConfigurationGroupService(mockRepo, NewTestTracer())

	name := "testGroup"
	version := model.Version{Major: 1, Minor: 0, Patch: 0}
	labels := "label1"
	configGroup := &model.ConfigurationGroup{
		Name:    name,
		Version: version,
	}

	mockRepo.On("GetGroupByParams", name, model.ToString(version), labels, mock.Anything).Return(configGroup, nil)

	retrievedGroup, err := service.Get(name, version, labels, context.Background())
	assert.NoError(t, err)
	assert.Equal(t, configGroup, retrievedGroup)

	mockRepo.AssertExpectations(t)
}

func TestConfigurationGroupService_Delete(t *testing.T) {
	mockRepo := new(repositories.MockConfigRepository)
	service := services.NewConfigurationGroupService(mockRepo, NewTestTracer())

	name := "testGroup"
	version := "1.0.0"
	labels := "label1"

	mockRepo.On("DeleteGroupByParams", name, version, labels, mock.Anything).Return(nil)

	err := service.Delete(name, version, labels, context.Background())
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
