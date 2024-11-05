package internal

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
)

var Server http.Server

func NewReverseProxy(conf *Config) error {
	slog.Debug("starting up reverse proxy", "port", conf.Port, "defaultTargetHost", conf.DefaultTargetHost, "services", len(conf.Services))
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, service := range conf.Services {
		serviceProxyTarget := proxyTargetDirector(service.TargetHost)
		for _, path := range service.Paths {
			mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
				slog.Info(fmt.Sprintf("[%s] %s%s => %s%s", r.Method, r.Host, r.RequestURI, service.TargetHost, r.RequestURI))
				serviceProxyTarget.ServeHTTP(w, r)
			})
		}
	}

	defaultTarget := proxyTargetDirector(conf.DefaultTargetHost)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info(fmt.Sprintf("[%s] %s%s => %s%s", r.Method, r.Host, r.RequestURI, conf.DefaultTargetHost, r.RequestURI))
		defaultTarget.ServeHTTP(w, r)
	})

	Server = http.Server{
		Addr:    fmt.Sprintf(":%s", conf.Port),
		Handler: mux,
	}
	err := Server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func proxyTargetDirector(host string) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		//req.Header.Add("X-Origin-Host", r.Host) // TODO: Needed?
		req.URL.Scheme = "http"
		req.URL.Host = host
		req.Host = host
	}
	return &httputil.ReverseProxy{Director: director}
}
