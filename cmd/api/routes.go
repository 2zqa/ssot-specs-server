package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	apiRootPath     = "/v1/"
	healthcheckPath = apiRootPath + "healthcheck"
	devicePath      = apiRootPath + "devices"
	devicePathUUID  = devicePath + "/:uuid"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Error handling
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Healthcheck
	router.HandlerFunc(http.MethodGet, healthcheckPath, app.healthcheckHandler)

	// Devices
	router.HandlerFunc(http.MethodGet, devicePath, app.requireOIDCAuthentication(app.listDevicesHandler))
	router.HandlerFunc(http.MethodPost, devicePath, app.requireAPIKeyAuthentication(app.createDeviceHandler))
	router.HandlerFunc(http.MethodPut, devicePathUUID, app.requireAPIKeyAuthentication(app.updateDeviceHandler))
	router.HandlerFunc(http.MethodGet, devicePathUUID, app.requireOIDCAuthentication(app.showDeviceHandler))
	router.HandlerFunc(http.MethodDelete, devicePathUUID, app.requireOIDCAuthentication(app.deleteDeviceHandler))

	// Middleware
	return app.recoverPanic(app.enableCORS(app.authenticateJwt(app.authenticateAPIKey(router))))
}
