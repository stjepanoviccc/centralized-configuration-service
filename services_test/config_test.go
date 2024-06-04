package services_test

import (
	"ars_projekat/model"
	"ars_projekat/repositories"
	"ars_projekat/services"
	"context"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func NewTestTracer() trace.Tracer {
	tp := sdktrace.NewTracerProvider()
	tracer := tp.Tracer("test-")
	return tracer
}
func TestConfigurationService_Add(t *testing.T) {
	config := &model.Configuration{Name: "testConfig", Version: model.Version{Major: 1, Minor: 0, Patch: 0}}

	mockRepo := new(repositories.MockConfigRepository)
	mockRepo.On("Add", config, mock.Anything).Return(config, nil)

	service := services.NewConfigurationService(mockRepo, NewTestTracer())

	err := service.Add(config, context.Background())
	assert.NoError(t, err)
	//assert.Equal(t, config, "")

	// Check if the configuration was added (using expectations)
	mockRepo.AssertCalled(t, "Add", config, mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestConfigurationService_Get(t *testing.T) {
	mockRepo := new(repositories.MockConfigRepository)
	service := services.NewConfigurationService(mockRepo, NewTestTracer())

	config := &model.Configuration{Name: "testConfig", Version: model.Version{Major: 1, Minor: 0, Patch: 0}}
	mockRepo.On("GetById", config.Name, model.ToString(config.Version), mock.Anything).Return(config, nil)

	retrievedConfig, err := service.Get(config.Name, model.ToString(config.Version), context.Background())
	assert.NoError(t, err)
	assert.Equal(t, config, retrievedConfig)

	mockRepo.AssertExpectations(t)
}

func TestConfigurationService_Delete(t *testing.T) {
	mockRepo := new(repositories.MockConfigRepository)
	service := services.NewConfigurationService(mockRepo, NewTestTracer())

	config := &model.Configuration{Name: "testConfig", Version: model.Version{Major: 1, Minor: 0, Patch: 0}}
	mockRepo.On("Delete", config.Name, model.ToString(config.Version), mock.Anything).Return(nil)

	err := service.Delete(*config, context.Background())
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
