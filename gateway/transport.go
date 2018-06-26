package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

//URL item in config YAML
type URL struct {
	Path   string `config:"path,required"`
	Method string `config:"method,required"`

	NSQ struct {
		Topic string `config:"topic"`
	}

	HTTP struct {
		Host string `config:"host"`
		Port string `config:"port"`
	}
}

type drivers struct {
	ctx     context.Context
	url     URL
	service Service
}

func (u *URL) buildHTTP() string {
	return fmt.Sprintf("http://%s:%s", u.HTTP.Host, u.HTTP.Port)
}

func (u *URL) isHTTP() bool {
	return u.HTTP.Host != ""
}

func (u *URL) isNsq() bool {
	return u.NSQ.Topic != ""
}

//MaketHTTPHandler creates handlers for Router
func MaketHTTPHandler(ctx context.Context, urls []URL, s Service, logger *logrus.Logger) http.Handler {
	r := mux.NewRouter()
	for _, u := range urls {
		switch {
		case u.isHTTP():
			director := func(req *http.Request) {
				out, err := url.Parse(u.buildHTTP())
				if err != nil {
					logger.Error(err)
				}
				req.URL.Scheme = out.Scheme
				req.URL.Host = out.Host
				req.URL.RawQuery = out.RawQuery
			}
			proxy := &httputil.ReverseProxy{Director: director}
			r.HandleFunc(u.Path, httpHandler(proxy)).Methods(u.Method)
		case u.isNsq():
			r.HandleFunc(u.Path, nsqHandler(ctx, u, s)).Methods(u.Method)
		}
	}
	r.Handle("/metrics", prometheus.Handler())
	return r
}

func httpHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
	}
}

func nsqHandler(ctx context.Context, u URL, s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			renderError(400, err, w)
			return
		}

		if len(data) == 0 {
			renderError(400, err, w)
			return
		}
		vars := mux.Vars(r)
		if err = s.Send(ctx, vars, data, u.NSQ.Topic); err != nil {
			renderError(500, err, w)
			return
		}
		return
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
