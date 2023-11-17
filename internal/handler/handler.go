package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thegoodparticle/vehicle-data-layer/internal/controller/vehicle"
	"github.com/thegoodparticle/vehicle-data-layer/internal/model"
	"github.com/thegoodparticle/vehicle-data-layer/internal/store"
	HttpStatus "github.com/thegoodparticle/vehicle-data-layer/internal/utils/http"
)

type Handler struct {
	Interface

	Controller vehicle.Interface
}

func NewHandler(repository store.Interface) Interface {
	return &Handler{
		Controller: vehicle.NewController(repository),
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if chi.URLParam(r, "RegID") != "" {
		h.getOne(w, r)
	} else {
		h.getAll(w, r)
	}
}

func (h *Handler) getOne(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "RegID")

	response, err := h.Controller.ListOne(ID)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, response)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	response, err := h.Controller.ListAll()
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, response)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	vehicleViolationsBody, err := h.getBodyAndValidate(r)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	ID, err := h.Controller.Create(vehicleViolationsBody)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, map[string]interface{}{"registration_id": ID})
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "RegID")

	vehicleViolationsBody, err := h.getBodyAndValidate(r)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	if err := h.Controller.Update(ID, vehicleViolationsBody); err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "RegID")

	if err := h.Controller.Remove(ID); err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusOK(w, r, "UP")
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) getBodyAndValidate(r *http.Request) (*model.VehicleViolations, error) {
	vehicleViolationsBody := &model.VehicleViolations{}
	body, err := model.ConvertIoReaderToStruct(r.Body, vehicleViolationsBody)
	if err != nil {
		return &model.VehicleViolations{}, errors.New("body is required")
	}

	vehicleBodyParsed, err := model.InterfaceToModel(body)
	if err != nil {
		return &model.VehicleViolations{}, errors.New("error on convert body to model")
	}

	if vehicleBodyParsed.VehicleRegID == "" {
		return &model.VehicleViolations{}, errors.New("registration ID is required")
	}

	log.Printf("successful parse of request body. %+v", vehicleBodyParsed)

	return vehicleBodyParsed, nil
}
