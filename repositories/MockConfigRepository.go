package repositories

import (
	"ars_projekat/model"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockConfigRepository struct {
	mock.Mock
}

func NewMockConfigRepository() *MockConfigRepository {
	return &MockConfigRepository{}
}

func (m *MockConfigRepository) GetAll(ctx context.Context) ([]model.Configuration, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Configuration), args.Error(1)
}

func (m *MockConfigRepository) GetById(name string, version string, ctx context.Context) (*model.Configuration, error) {
	args := m.Called(name, version, ctx)
	return args.Get(0).(*model.Configuration), args.Error(1)
}

func (m *MockConfigRepository) Delete(name string, version string, ctx context.Context) error {
	args := m.Called(name, version, ctx)
	return args.Error(0)
}

func (m *MockConfigRepository) Add(config *model.Configuration, ctx context.Context) (*model.Configuration, error) {
	args := m.Called(config, ctx)
	return args.Get(0).(*model.Configuration), args.Error(1)
}

func (m *MockConfigRepository) GetAllGroups(ctx context.Context) ([]model.ConfigurationGroup, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.ConfigurationGroup), args.Error(1)
}

func (m *MockConfigRepository) GetGroupByParams(name string, version string, labels string, ctx context.Context) (*model.ConfigurationGroup, error) {
	args := m.Called(name, version, labels)
	return args.Get(0).(*model.ConfigurationGroup), args.Error(1)
}

func (m *MockConfigRepository) AddGroup(name string, version string, labels string, configs model.Configuration, ctx context.Context) error {
	args := m.Called(name, version, labels, configs, ctx)
	return args.Error(0)
}

func (m *MockConfigRepository) DeleteGroupById(name string, version string, ctx context.Context) error {
	args := m.Called(name, version, ctx)
	return args.Error(0)
}

func (m *MockConfigRepository) DeleteGroupByParams(name string, version string, labels string, ctx context.Context) error {
	args := m.Called(name, version, labels, ctx)
	return args.Error(0)
}

func (m *MockConfigRepository) GetIdempotencyRequestByKey(key string, ctx context.Context) (bool, error) {
	args := m.Called(key, ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockConfigRepository) AddIdempotencyRequest(req *model.IdempotencyRequest, ctx context.Context) (*model.IdempotencyRequest, error) {
	args := m.Called(req, ctx)
	return args.Get(0).(*model.IdempotencyRequest), args.Error(1)
}
