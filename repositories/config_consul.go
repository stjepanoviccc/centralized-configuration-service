package repositories

import (
	"ars_projekat/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"log"
	"os"

	"github.com/hashicorp/consul/api"

	"go.opentelemetry.io/otel/codes"
)

type ConfigRepository struct {
	cli    *api.Client
	logger *log.Logger
	Tracer trace.Tracer
}

// Config
func New(logger *log.Logger, tracer trace.Tracer) (*ConfigRepository, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigRepository{
		cli:    client,
		logger: logger,
		Tracer: tracer,
	}, nil
}

func (cr *ConfigRepository) GetAll(ctx context.Context) ([]model.Configuration, error) {
	_, span := cr.Tracer.Start(ctx, "ConfigRepository.GetAll")
	defer span.End()
	kv := cr.cli.KV()
	data, _, err := kv.List(allConfigs, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var configurations []model.Configuration
	for _, pair := range data {
		configuration := &model.Configuration{}
		err = json.Unmarshal(pair.Value, configuration)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		configurations = append(configurations, *configuration)
	}

	span.SetStatus(codes.Ok, "Success fetching configurations")
	return configurations, nil
}

func (cr *ConfigRepository) GetById(name string, version string, ctx context.Context) (*model.Configuration, error) {
	_, span := cr.Tracer.Start(ctx, "ConfigRepository.GetById")
	defer span.End()

	kv := cr.cli.KV()
	data, _, err := kv.Get(ConstructConfigKey(name, version), nil)
	if data == nil {
		return nil, errors.New("not found")
	}
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	configuration := &model.Configuration{}
	err = json.Unmarshal(data.Value, configuration)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "Success fetching configuration")
	return configuration, nil
}

func (cr *ConfigRepository) Delete(name string, version string, ctx context.Context) error {
	_, span := cr.Tracer.Start(ctx, "ConfigRepository.Delete")
	defer span.End()

	kv := cr.cli.KV()

	_, err := kv.Delete(ConstructConfigKey(name, version), nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "Success deleting configuration")
	return nil
}

func (cr *ConfigRepository) Add(config *model.Configuration, ctx context.Context) (*model.Configuration, error) {
	_, span := cr.Tracer.Start(ctx, "ConfigRepository.Add")
	defer span.End()

	kv := cr.cli.KV()
	version := model.ToString(config.Version)

	data, err := json.Marshal(config)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	keyValue := &api.KVPair{Key: ConstructConfigKey(config.Name, version), Value: data}
	_, err = kv.Put(keyValue, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "Successfully added Configuration")
	return config, nil
}

// Config group
func (cr *ConfigRepository) GetAllGroups(ctx context.Context) ([]model.ConfigurationGroup, error) {
	_, span := cr.Tracer.Start(ctx, "ConfigGroupRepository.GetAllGroups")
	defer span.End()

	kv := cr.cli.KV()
	data, _, err := kv.List(allGroups, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var groups []model.ConfigurationGroup
	for _, pair := range data {
		cg := &model.ConfigurationGroup{}
		err = json.Unmarshal(pair.Value, cg)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		groups = append(groups, *cg)
	}

	span.SetStatus(codes.Ok, "Success fetching all config groups")
	return groups, nil
}

func (cr *ConfigRepository) GetGroupByParams(name string, version string, labels string, ctx context.Context) (*model.ConfigurationGroup, error) {
	_, span := cr.Tracer.Start(ctx, "ConfigGroupRepository.GetGroupByParams")
	defer span.End()

	kv := cr.cli.KV()

	var key string
	if len(labels) == 0 {
		key = ConstructConfigGroupKey(name, version, "", "")
	} else {
		key = ConstructConfigGroupKey(name, version, labels, "")
	}

	data, _, err := kv.List(key, nil)
	if data == nil {
		return nil, err
	}
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	cg := &model.ConfigurationGroup{}
	cg.Name = name
	ver, err := model.ToVersion(version)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	cg.Version = *ver

	for _, pair := range data {
		config := &model.Configuration{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		cg.Configurations = append(cg.Configurations, *config)
	}

	span.SetStatus(codes.Ok, "Success fetching group by parameters")
	return cg, nil
}

func (cr *ConfigRepository) AddGroup(name string, version string, labels string, configs model.Configuration, ctx context.Context) error {
	_, span := cr.Tracer.Start(ctx, "ConfigGroupRepository.AddGroup")
	defer span.End()

	kv := cr.cli.KV()

	data, err := json.Marshal(configs)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	keyValue := &api.KVPair{Key: ConstructConfigGroupKey(name, version, labels, configs.Name), Value: data}
	_, err = kv.Put(keyValue, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "Successfully added group")
	return nil
}

func (cr *ConfigRepository) DeleteGroupById(name string, version string, ctx context.Context) error {
	_, span := cr.Tracer.Start(ctx, "ConfigGroupRepository.DeleteGroupById")
	defer span.End()

	kv := cr.cli.KV()

	_, err := kv.Delete(ConstructConfigKey(name, version), nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "Successfully deleted configuration group")
	return nil
}

func (cr *ConfigRepository) DeleteGroupByParams(name string, version string, labels string, ctx context.Context) error {
	_, span := cr.Tracer.Start(ctx, "ConfigGroupRepository.DeleteGroupByParams")
	defer span.End()

	kv := cr.cli.KV()

	var key string
	if len(labels) == 0 {
		key = ConstructConfigGroupKey(name, version, "", "")
	} else {
		key = ConstructConfigGroupKey(name, version, labels, "")
	}
	_, err := kv.DeleteTree(key, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "Successfully deleted configuration group")
	return nil
}

func (cr *ConfigRepository) GetIdempotencyRequestByKey(key string, ctx context.Context) (bool, error) {
	_, span := cr.Tracer.Start(ctx, "Repository.IdempotencyRequest")
	defer span.End()

	kv := cr.cli.KV()

	data, _, err := kv.Get(ConstructIdempotencyRequestKey(key), nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return false, err
	}
	if data == nil {
		return false, nil
	}

	span.SetStatus(codes.Ok, "Success, finishing up")
	return true, nil
}

func (cr *ConfigRepository) AddIdempotencyRequest(req *model.IdempotencyRequest, ctx context.Context) (*model.IdempotencyRequest, error) {
	_, span := cr.Tracer.Start(ctx, "Repository.AddIdempotencyRequest")
	defer span.End()

	kv := cr.cli.KV()

	data, err := json.Marshal(req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	keyValue := &api.KVPair{Key: ConstructIdempotencyRequestKey(req.Key), Value: data}
	_, err = kv.Put(keyValue, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "Success, finishing up")
	return req, nil
}

type IConfigRepository interface {
	GetAll(ctx context.Context) ([]model.Configuration, error)
	GetById(name string, version string, ctx context.Context) (*model.Configuration, error)
	Delete(name string, version string, ctx context.Context) error
	Add(config *model.Configuration, ctx context.Context) (*model.Configuration, error)
	GetAllGroups(ctx context.Context) ([]model.ConfigurationGroup, error)
	GetGroupByParams(name string, version string, labels string, ctx context.Context) (*model.ConfigurationGroup, error)
	AddGroup(name string, version string, labels string, configs model.Configuration, ctx context.Context) error
	DeleteGroupById(name string, version string, ctx context.Context) error
	DeleteGroupByParams(name string, version string, labels string, ctx context.Context) error
	GetIdempotencyRequestByKey(key string, ctx context.Context) (bool, error)
	AddIdempotencyRequest(req *model.IdempotencyRequest, ctx context.Context) (*model.IdempotencyRequest, error)
}
