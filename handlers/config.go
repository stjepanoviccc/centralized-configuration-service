package handlers

import (
	"ars_projekat/model"
	"ars_projekat/services"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ConfigurationHandler struct {
	Tracer  trace.Tracer
	Service services.ConfigurationService
}

func NewConfigurationHandler(service services.ConfigurationService, tracer trace.Tracer) ConfigurationHandler {
	return ConfigurationHandler{
		Service: service,
		Tracer:  tracer,
	}
}

// swagger:route GET /configs/{name}/{version} configuration getConfiguration
// Get configuration by name and version
//
// responses:
//
//	404: ErrorResponse
//	200: Configuration
func (c ConfigurationHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, span := c.Tracer.Start(r.Context(), "ConfigurationHandler.Get")
	defer span.End()

	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	config, err := c.Service.Get(name, version, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(ctx, w, config, http.StatusOK)
	span.SetStatus(codes.Ok, "")
}

// swagger:route POST /configs configuration upsertConfiguration
// Add or update a configuration
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	409: ErrorResponse
//	201: Configuration
func (c ConfigurationHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	ctx, span := c.Tracer.Start(r.Context(), "ConfigurationHandler.Upsert")
	defer span.End()

	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediaType != "application/json" {
		err := errors.New("expect application/json Content-Type")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	cfg, err := decodeBody(r.Body)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ver := model.ToString(cfg.Version)
	check, err := c.Service.Get(cfg.Name, ver, ctx)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if check != nil {
		err := errors.New("config already exists")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	err = c.Service.Add(cfg, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderJSON(ctx, w, cfg, http.StatusCreated)
	span.SetStatus(codes.Ok, "")
}

// swagger:route DELETE /configs/{name}/{version} configuration deleteConfiguration
// Delete a configuration by name and version
//
// responses:
//
//	404: ErrorResponse
//	204: NoContent
func (c ConfigurationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := c.Tracer.Start(r.Context(), "ConfigurationHandler.Delete")
	defer span.End()

	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]

	config, err := c.Service.Get(name, version, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok := c.Service.Delete(*config, ctx)
	if ok != nil {
		span.SetStatus(codes.Error, ok.Error())
		http.Error(w, ok.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	span.SetStatus(codes.Ok, "")
}

func decodeBody(r io.Reader) (*model.Configuration, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var configuration model.Configuration
	if err := dec.Decode(&configuration); err != nil {
		return nil, err
	}
	return &configuration, nil
}

func renderJSON(ctx context.Context, w http.ResponseWriter, v interface{}, statusCode int) {
	marshal, err := json.Marshal(v)
	if err != nil {
		span := trace.SpanFromContext(ctx)
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err = w.Write(marshal); err != nil {
		span := trace.SpanFromContext(ctx)
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
