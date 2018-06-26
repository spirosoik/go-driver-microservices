package zombie

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

//MakeHTTPHandler creates an http handler
func MakeHTTPHandler(s Service, logger *logrus.Logger) http.Handler {
	r := mux.NewRouter()
	r.Path("/drivers/{id}").
		Methods("GET").
		HandlerFunc(driverHandler(s, logger))

	r.Handle("/metrics", prometheus.Handler())
	return r
}

func driverHandler(s Service, logger *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		d, err := s.IsZombie(vars["id"], 5)
		if err != nil {
			renderError(500, err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&d); err != nil {
			renderError(500, err, w)
		}
	}
}

func renderError(code int, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": code,
		"error":  err.Error(),
	})
}
