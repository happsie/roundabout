package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
)

func NewReverseProxy(conf *Config) error {
	slog.Debug("starting up reverse proxy", "port", conf.Port, "defaultTargetHost", conf.DefaultTargetHost, "services", len(conf.Services))
	mux := http.NewServeMux()

	defaultTarget := defaultProxyTarget(conf)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info(fmt.Sprintf("[%s] %s%s => %s%s", r.Method, r.Host, r.RequestURI, conf.DefaultTargetHost, r.RequestURI))
		defaultTarget.ServeHTTP(w, r)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", conf.Port), mux)
	if err != nil {
		return err
	}
	return nil
}

func defaultProxyTarget(conf *Config) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		//req.Header.Add("X-Origin-Host", r.Host) // TODO: Needed?
		req.URL.Scheme = "http"
		req.URL.Host = conf.DefaultTargetHost
		req.Host = conf.DefaultTargetHost
	}
	return &httputil.ReverseProxy{Director: director}
}
