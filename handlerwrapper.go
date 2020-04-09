package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
	"net/http/httputil"
)

type HandlerWrapper struct {
	Usecase UsecaseInterface
	Configs *Config
}

func NewHandlerWrapper() (*HandlerWrapper, error) {
	c, err := NewConfig()
	if err != nil {
		return nil, err
	}

	usecase, err := NewUsecase(c)
	if err != nil {
		return nil, err
	}

	return &HandlerWrapper{Usecase: usecase, Configs: c}, nil
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func respondError(w http.ResponseWriter, code int, err error) {
	respondJSON(w, code, map[string]string{"error": err.Error()})
}

func (a *HandlerWrapper) WrapHandler(h http.Handler) http.Handler {

	authFunc := func(w http.ResponseWriter, r *http.Request) {

		dumpRequest(r)
		if r.Method == "OPTIONS" {
			respondJSON(w, http.StatusOK, "")
			return
		}

		h.ServeHTTP(w, r)
	}
	return allowCORS(http.HandlerFunc(authFunc))

}

func dumpRequest(r *http.Request) {

	dump, _ := httputil.DumpRequest(r, true)
	log.Info(string(dump))
}

func allowCORS(h http.Handler) http.Handler {
	options := []handlers.CORSOption{
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "ClientId", "Env", "Access-Control-Allow-Origin"}),
	}

	return handlers.CORS(options...)(h)
}
