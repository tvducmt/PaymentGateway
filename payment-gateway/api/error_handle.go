package api

import (
	"fmt"
	"net/http"

	raven "github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
)

// ErrorInterface interface for StatusError
type ErrorInterface interface {
	error
	Status() int
	Data() []byte
}

// StatusError a custom error with http status code
type StatusError struct {
	Code    int
	Err     error
	Results []byte
}

// Error for satisfy error interface
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Status return http status code
func (se StatusError) Status() int {
	return se.Code
}

// Data return http result
func (se StatusError) Data() []byte {
	return se.Results
}

type responseError struct {
	Message string
}

// HandleAPI get handler with custom return error
type HandleAPI func(w http.ResponseWriter, r *http.Request) error

// Reg return http.HandleFunc satisfy interface
func (h HandleAPI) Reg(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		switch e := err.(type) {
		case ErrorInterface:
			switch e.Status() {
			case 500:
				//log.Printf("HTTP %d - %s", e.Status(), e)
				log.WithFields(log.Fields{
					"Status": e.Status(),
					"Error":  e.Error(),
				}).Info("System Error")

				raven.CaptureErrorAndWait(e, nil)
				w.Write([]byte(`{"error": "Something went wrong. please try again"}`))
			case 501:
				log.WithFields(log.Fields{
					"Status": e.Status(),
					"Error":  e.Error(),
				}).Info("System Error")

				raven.CaptureErrorAndWait(e, nil)
				w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, e.Error())))
			case 401:
				log.WithFields(log.Fields{
					"Status": e.Status(),
					"Error":  e.Error(),
				}).Info("Something went wrong, please try again")
				raven.CaptureErrorAndWait(e, nil)
				w.Write([]byte(`{"error": "Something went wrong. please try again"}`))
			case 200:
				w.Write(e.Data())
			case 201:
				raven.CaptureErrorAndWait(e, nil)
				w.Write(e.Data())
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
