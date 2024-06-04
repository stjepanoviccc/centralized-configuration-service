package services

import (
	"ars_projekat/model"
	"ars_projekat/repositories"
	"context"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ConfigurationGroupService struct {
	repo   repositories.IConfigRepository
	Tracer trace.Tracer
}

func NewConfigurationGroupService(repo repositories.IConfigRepository, tracer trace.Tracer) ConfigurationGroupService {
	return ConfigurationGroupService{
		repo:   repo,
		Tracer: tracer,
	}
}

func (s ConfigurationGroupService) Add(configGroup model.ConfigurationGroup, ctx context.Context) error {
	ctx, span := s.Tracer.Start(ctx, "ConfigurationGroupService.Add")
	defer span.End()

	name := configGroup.Name
	version := model.ToString(configGroup.Version)
	for _, v := range configGroup.Configurations {
		labels := model.SortLabels(v.Labels)
		err := s.repo.AddGroup(name, version, labels, v, ctx)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return err
		}
	}

	span.SetStatus(codes.Ok, "SERVICE - Success")
	return nil
}

func (s ConfigurationGroupService) Save(configGroup *model.ConfigurationGroup, ctx context.Context) error {
	ctx, span := s.Tracer.Start(ctx, "ConfigurationGroupService.Save")
	defer span.End()

	for _, v := range configGroup.Configurations {
		labels := model.SortLabels(v.Labels)
		err := s.repo.AddGroup(configGroup.Name, model.ToString(configGroup.Version), labels, v, ctx)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return err
		}
	}

	span.SetStatus(codes.Ok, "SERVICE - Success")
	return nil
}

func (s ConfigurationGroupService) Get(name string, version model.Version, labels string, ctx context.Context) (*model.ConfigurationGroup, error) {
	ctx, span := s.Tracer.Start(ctx, "ConfigurationGroupService.Get")
	defer span.End()

	return s.repo.GetGroupByParams(name, model.ToString(version), labels, ctx)
}

func (s ConfigurationGroupService) Delete(name string, version string, labels string, ctx context.Context) error {
	ctx, span := s.Tracer.Start(ctx, "ConfigurationGroupService.Get")
	defer span.End()

	return s.repo.DeleteGroupByParams(name, version, labels, ctx)
}
