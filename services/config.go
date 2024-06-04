package services

import (
	"ars_projekat/model"
	"ars_projekat/repositories"
	"context"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ConfigurationService struct {
	repo   repositories.IConfigRepository
	Tracer trace.Tracer
}

func NewConfigurationService(repo repositories.IConfigRepository, tracer trace.Tracer) ConfigurationService {
	return ConfigurationService{
		repo:   repo,
		Tracer: tracer,
	}
}

func (s ConfigurationService) Add(config *model.Configuration, ctx context.Context) error {
	ctx, span := s.Tracer.Start(ctx, "ConfigurationService.Add")
	defer span.End()

	_, err := s.repo.Add(config, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "SERVICE - Success")
	return nil
}

func (s ConfigurationService) Get(name string, version string, ctx context.Context) (*model.Configuration, error) {
	ctx, span := s.Tracer.Start(ctx, "ConfigurationService.Get")
	defer span.End()

	return s.repo.GetById(name, version, ctx)
}

func (s ConfigurationService) Delete(config model.Configuration, ctx context.Context) error {
	ctx, span := s.Tracer.Start(ctx, "ConfigurationService.Delete")
	defer span.End()

	ver := model.ToString(config.Version)

	return s.repo.Delete(config.Name, ver, ctx)
}
