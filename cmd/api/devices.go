package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/2zqa/ssot-specs-server/internal/data"
	"github.com/2zqa/ssot-specs-server/internal/validator"
	"github.com/google/uuid"
)

func (app *application) createDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID    uuid.UUID  `json:"uuid"`
		Specs data.Specs `json:"specs"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	device := &data.Device{
		ID:    input.ID,
		Specs: input.Specs,
	}

	v := validator.New()

	if data.ValidateDevice(v, device); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Devices.Insert(device, app.models.Search)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateDeviceID):
			v.AddError("device", "Device already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/devices/%s", device.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"device": device}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Specs data.Specs `json:"specs"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	device := &data.Device{
		ID:    uuid,
		Specs: input.Specs,
	}

	v := validator.New()

	if data.ValidateDevice(v, device); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	newDeviceCreated, err := app.updateOrInsertDevice(device)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	statusCode := http.StatusOK
	if newDeviceCreated {
		statusCode = http.StatusCreated
	}

	err = app.writeJSON(w, statusCode, envelope{"device": device}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateOrInsertDevice(device *data.Device) (newDeviceCreated bool, err error) {
	_, err = app.models.Devices.Get(device.ID)

	switch err {
	case data.ErrRecordNotFound:
		if err = app.models.Devices.Insert(device, app.models.Search); err != nil {
			return false, err
		}
		return true, nil
	case nil:
		if err = app.models.Devices.Update(device, app.models.Search); err != nil {
			return false, err
		}
		return false, nil
	default:
		return false, err
	}
}

func (app *application) deleteDeviceHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Devices.Delete(uuid, app.models.Search)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showDeviceHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	device, err := app.models.Devices.Get(uuid)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"device": device}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listDevicesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SearchTerms string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	input.SearchTerms = app.readString(qs, "q", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "uuid")
	input.Filters.SortSafelist = []string{"updated_at", "-updated_at"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	devices, metadata, err := app.models.Devices.GetAll(input.SearchTerms, input.Filters)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"devices": devices, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
