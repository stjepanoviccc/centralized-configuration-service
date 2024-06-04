package services

import (
	"ars_projekat/model"
	"ars_projekat/repositories"
	"context"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type IdempotencyService struct {
	repo   repositories.ConfigRepository
	Tracer trace.Tracer
}

func NewIdempotencyService(repo repositories.ConfigRepository, tracer trace.Tracer) IdempotencyService {
	return IdempotencyService{
		repo:   repo,
		Tracer: tracer,
	}
}

func (i IdempotencyService) Add(req *model.IdempotencyRequest, ctx context.Context) error {
	ctx, span := i.Tracer.Start(ctx, "IdempotencyService.Add")
	defer span.End()

	_, err := i.repo.AddIdempotencyRequest(req, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "SERVICE - Success")
	return nil
}

func (i IdempotencyService) Get(key string, ctx context.Context) (bool, error) {
	ctx, span := i.Tracer.Start(ctx, "IdempotencyService.Get")
	defer span.End()

	exists, err := i.repo.GetIdempotencyRequestByKey(key, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return false, err
	}

	span.SetStatus(codes.Ok, "SERVICE - Success")
	return exists, nil
}
