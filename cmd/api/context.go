package main

import (
	"context"
	"fmt"
	"net/http"
)

type contextKey string

const oidcContextKey = contextKey("oidc-authenticated")
const apiKeyContextKey = contextKey("api-key-authenticated")
const missingKeyMessage = "missing %s value in request context"

func (app *application) contextSetIsOIDCAuthenticated(r *http.Request, isAuth bool) *http.Request {
	ctx := context.WithValue(r.Context(), oidcContextKey, isAuth)
	return r.WithContext(ctx)
}

func (app *application) contextIsOIDCAuthenticated(r *http.Request) bool {
	isAuth, ok := r.Context().Value(oidcContextKey).(bool)
	if !ok {
		panic(fmt.Sprintf(missingKeyMessage, oidcContextKey))
	}

	return isAuth
}

func (app *application) contextSetIsAPIKeyAuthenticated(r *http.Request, isAuth bool) *http.Request {
	ctx := context.WithValue(r.Context(), apiKeyContextKey, isAuth)
	return r.WithContext(ctx)
}

func (app *application) contextIsAPIKeyAuthenticated(r *http.Request) bool {
	isAuth, ok := r.Context().Value(apiKeyContextKey).(bool)
	if !ok {
		panic(fmt.Sprintf(missingKeyMessage, apiKeyContextKey))
	}

	return isAuth
}
