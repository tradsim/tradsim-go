package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/adaptlog"
)

type statusLoggingResponseWriter struct {
	status              int
	statusHeaderWritten bool
	w                   http.ResponseWriter
}

func (w *statusLoggingResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *statusLoggingResponseWriter) Write(d []byte) (int, error) {

	value, err := w.w.Write(d)
	if err != nil {
		return value, err
	}

	if !w.statusHeaderWritten {
		w.status = http.StatusOK
		w.statusHeaderWritten = true
	}

	return value, err
}

func (w *statusLoggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.w.WriteHeader(code)
	w.statusHeaderWritten = true
}

// DefaultMiddleware which handles Logging and Recover middleware
func DefaultMiddleware(next httprouter.Handle) httprouter.Handle {
	return LoggingMiddleware(RecoveryMiddleware(next))
}

// DefaultGETValidationMiddleware which handles the POST and JSON validation, Logging and Recover middleware
func DefaultGETValidationMiddleware(next httprouter.Handle) httprouter.Handle {
	return GETValidationMiddleware(DefaultMiddleware(next))
}

// DefaultPOSTJSONValidationMiddleware which handles the POST and JSON validation, Logging and Recover middleware
func DefaultPOSTJSONValidationMiddleware(next httprouter.Handle) httprouter.Handle {
	return POSTJSONValidationMiddleware(DefaultMiddleware(next))
}

// DefaultPUTJSONValidationMiddleware which handles the POST and JSON validation, Logging and Recover middleware
func DefaultPUTJSONValidationMiddleware(next httprouter.Handle) httprouter.Handle {
	return PUTJSONValidationMiddleware(DefaultMiddleware(next))
}

// DefaultDELETEValidationMiddleware which handles the POST and JSON validation, Logging and Recover middleware
func DefaultDELETEValidationMiddleware(next httprouter.Handle) httprouter.Handle {
	return DELETEValidationMiddleware(DefaultMiddleware(next))
}

// LoggingMiddleware for recovering from failed requests
func LoggingMiddleware(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		lw := &statusLoggingResponseWriter{-1, false, w}
		startTime := time.Now()
		next(lw, r, ps)
		adaptlog.NewStdLevelLogger("LoggingMiddleware").Infof("host=%s method=%s route=%s status=%d time=%s params=%s", r.Host, r.Method, r.URL.String(), lw.status, time.Since(startTime), ps)
	}
}

// RecoveryMiddleware for recovering from failed requests
func RecoveryMiddleware(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		defer func() {
			if err := recover(); err != nil {
				adaptlog.NewStdLevelLogger("RecoveryMiddleware").Errorf("[ERROR] %s", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next(w, r, ps)
	}
}

// POSTJSONValidationMiddleware validates incomming requests for POST method and JSON headers
func POSTJSONValidationMiddleware(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		logger := adaptlog.NewStdLevelLogger("POSTJSONValidationMiddleware")

		if r.Method != http.MethodPost {

			logger.Warnf("Http method POST was expected, but received %s instead", r.Method)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		contentType := r.Header.Get("Content-Type")

		if !strings.HasPrefix(contentType, "application/json") {

			logger.Warnf("Content type is not 'application/json', but %s instead", contentType)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		next(w, r, ps)
	}
}

// GETValidationMiddleware validates incoming requests for GET method
func GETValidationMiddleware(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		logger := adaptlog.NewStdLevelLogger("GETValidationMiddleware")

		if r.Method != http.MethodGet {

			logger.Warnf("Http method GET was expected, but received %s instead", r.Method)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		next(w, r, ps)
	}
}

// PUTJSONValidationMiddleware validates incomming requests for POST method and JSON headers
func PUTJSONValidationMiddleware(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		logger := adaptlog.NewStdLevelLogger("PUTJSONValidationMiddleware")

		if r.Method != http.MethodPut {

			logger.Warnf("Http method PUT was expected, but received %s instead", r.Method)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		contentType := r.Header.Get("Content-Type")

		if !strings.HasPrefix(contentType, "application/json") {

			logger.Warnf("Content type is not 'application/json', but %s instead", contentType)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		next(w, r, ps)
	}
}

// DELETEValidationMiddleware validates incoming requests for GET method
func DELETEValidationMiddleware(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		logger := adaptlog.NewStdLevelLogger("DELETEValidationMiddleware")

		if r.Method != http.MethodDelete {

			logger.Warnf("Http method DELETE was expected, but received %s instead", r.Method)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		next(w, r, ps)
	}
}
