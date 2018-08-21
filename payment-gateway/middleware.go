package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	raven "github.com/getsentry/raven-go"
	"github.com/go-chi/chi/middleware"
	"gitlab.com/rockship/payment-gateway/api"
)

// TrustDeviceTokenMiddleware ...
func TrustDeviceTokenMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != os.Getenv("TRUST_DEVICE_TOKEN") {
			http.Error(w, http.StatusText(404), http.StatusNotAcceptable)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// SentryLoggingMiddleware ...
func SentryLoggingMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					raven.CapturePanic(func() {
						log.Panic(rvr)
					}, nil)
				} else {
					fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
					debug.PrintStack()
				}
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// OauthAuthorizationMiddleware ...
func OauthAuthorizationMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		accessToken := ""
		bearer := r.Header.Get("Authorization")
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			accessToken = bearer[7:]
		}

		resp, user, err := api.CallApiVerifyAccessToken(accessToken)

		if resp.StatusCode < 200 && resp.StatusCode >= 300 || accessToken == "" || err != nil {
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user", *user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// JSONReturnMiddleware ...
func JSONReturnMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
