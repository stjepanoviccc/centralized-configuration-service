package handlers

import (
	"ars_projekat/model"
	"ars_projekat/services"
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

type ConfigurationGroupHandler struct {
	Tracer       trace.Tracer
	GroupService services.ConfigurationGroupService
}

func NewConfigurationGroupHandler(groupService services.ConfigurationGroupService, tracer trace.Tracer) ConfigurationGroupHandler {
	return ConfigurationGroupHandler{
		GroupService: groupService,
		Tracer:       tracer,
	}
}

// swagger:route GET /config-groups/{name}/{version}/{labels} configurationgroup getConfigurationGroup
// Get configuration group by name, version, and labels
//
// responses:
//
//	404: ErrorResponse
//	200: ConfigurationGroup
func (cg ConfigurationGroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, span := cg.Tracer.Start(r.Context(), "ConfigurationGroupHandler.Get")
	defer span.End()

	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	labels := strings.Split(mux.Vars(r)["labels"], ";")

	var labelString string
	for i, v := range labels {
		if i == len(labels)-1 {
			labelString += v
		} else {
			labelString += v
			labelString += ";"
		}
	}
	versionModel, err := model.ToVersion(version)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cGroup, err := cg.GroupService.Get(name, *versionModel, labelString, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if cGroup == nil {
		http.Error(w, "no content", http.StatusNoContent)
		return
	}

	renderJSON(ctx, w, cGroup)
	span.SetStatus(codes.Ok, "")
}

// swagger:route POST /config-groups/{name}/{version} configurationgroup addConfigurationToGroup
// Add configuration to a configuration group
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	409: ErrorResponse
//	201: ConfigurationGroup
func (cg ConfigurationGroupHandler) AddConfig(w http.ResponseWriter, r *http.Request) {
	ctx, span := cg.Tracer.Start(r.Context(), "ConfigurationGroupHandler.AddConfig")
	defer span.End()

	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionModel, err := model.ToVersion(version)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cGroup, err := cg.GroupService.Get(name, *versionModel, "", ctx)
	if cGroup == nil {
		err = errors.New("config not found")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	cType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(cType)
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

	config, err := decodeBody(r.Body)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, v := range cGroup.Configurations {
		if v.Name == config.Name && v.Version == config.Version {
			err := errors.New("config is already added")
			span.SetStatus(codes.Error, err.Error())
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
	}

	cGroup.Configurations = append(cGroup.Configurations, *config)
	err = cg.GroupService.Save(cGroup, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(ctx, w, cGroup)
	span.SetStatus(codes.Ok, "")
}

// swagger:route POST /config-groups configurationgroup upsertConfigurationGroup
// Add or update a configuration group
//
// responses:
//
//	415: ErrorResponse
//	400: ErrorResponse
//	201: ConfigurationGroup
func (cg ConfigurationGroupHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	ctx, span := cg.Tracer.Start(r.Context(), "ConfigurationGroupHandler.Upsert")
	defer span.End()

	cType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(cType)
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

	cfgGroup, err := decodeGroupBody(r.Body)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = cg.GroupService.Add(*cfgGroup, ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(ctx, w, cfgGroup)
	span.SetStatus(codes.Ok, "")
}

// swagger:route DELETE /config-groups/{name}/{version}/{labels} configurationgroup deleteConfigurationGroup
// Delete a configuration group by name, version, and labels
//
// responses:
//
//	404: ErrorResponse
//	204: NoContent
func (cg ConfigurationGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := cg.Tracer.Start(r.Context(), "ConfigurationGroupHandler.Delete")
	defer span.End()

	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	labels := strings.Split(mux.Vars(r)["labels"], ";")

	var labelString string
	for i, v := range labels {
		if i == len(labels)-1 {
			labelString += v
		} else {
			labelString += v
			labelString += ";"
		}
	}
	versionModel, err := model.ToVersion(version)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	check, err := cg.GroupService.Get(name, *versionModel, labelString, ctx)
	if check == nil {
		err = errors.New("config not found")
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok := cg.GroupService.Delete(name, version, labelString, ctx)
	if ok != nil {
		span.SetStatus(codes.Error, ok.Error())
		http.Error(w, ok.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	span.SetStatus(codes.Ok, "")
}

func decodeGroupBody(r io.Reader) (*model.ConfigurationGroup, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var cg model.ConfigurationGroup
	if err := dec.Decode(&cg); err != nil {
		return nil, err
	}

	return &cg, nil
}
