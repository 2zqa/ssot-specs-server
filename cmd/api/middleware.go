package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or
			// not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the
				// response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been
				// sent.
				w.Header().Set("Connection", "close")
				// log the error and send the client a 500 Internal Server Error response.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticateAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip if authentication is disabled
		if app.config.disableAuth {
			r = app.contextSetIsAPIKeyAuthenticated(r, true)
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Add("Vary", "X-API-KEY")

		apiKey := r.Header.Get("X-API-KEY")
		if apiKey == "" {
			r = app.contextSetIsAPIKeyAuthenticated(r, false)
			next.ServeHTTP(w, r)
			return
		}

		if apiKey != app.config.apiKey {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		app.logger.Info("Authenticated through API key")
		r = app.contextSetIsAPIKeyAuthenticated(r, true)
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticateJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip if authentication is disabled
		if app.config.disableAuth {
			r = app.contextSetIsOIDCAuthenticated(r, true)
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Add("Vary", "Authorization")

		// Continue if no authorization is provided.
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = app.contextSetIsOIDCAuthenticated(r, false)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		rawIDToken := headerParts[1]

		ctx := context.Background()
		// verifier.Verify checks issuer, audience, signatures and expiry time.
		idToken, err := app.oidcInfo.verifier.Verify(ctx, rawIDToken)
		if err != nil {
			app.logger.Info("Authentication failed", "err", err)
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Extract custom claims
		var claims struct {
			Verified bool   `json:"email_verified"`
			Name     string `json:"name"`
		}
		if err := idToken.Claims(&claims); err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Check if users e-mail is verified
		if !claims.Verified {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Call the next handler in the chain.
		app.logger.Info("Authenticated through OIDC", "claims", claims)
		r = app.contextSetIsOIDCAuthenticated(r, true)
		next.ServeHTTP(w, r)
	})
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		for i := range app.config.cors.trustedOrigins {
			if origin != app.config.cors.trustedOrigins[i] {
				continue
			}

			// Origin is on allowed list, set the Access-Control-Allow-Origin
			w.Header().Set("Access-Control-Allow-Origin", origin)

			if !isPreflightRequest(r.Method, r.Header) {
				break
			}

			// Preflight request, so specify allowed methods and headers as well
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, X-API-KEY, Content-Type")

			// The use of StatusOK over StatusNoContent is for compatibility reasons: https://stackoverflow.com/a/58794243
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireOIDCAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuthenticated := app.contextIsOIDCAuthenticated(r)

		if !isAuthenticated {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAPIKeyAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuthenticated := app.contextIsAPIKeyAuthenticated(r)

		if !isAuthenticated {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isPreflightRequest(method string, header http.Header) bool {
	return method == http.MethodOptions && header.Get("Access-Control-Request-Method") != ""
}
