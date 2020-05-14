package http

import (
	"github.com/alabianca/kadbox/log"
	"net/http"
)

type AppHandler struct {
	StorageHandler http.Handler
	PingHandler   http.Handler
}

func (a *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	path := r.URL.Path

	log.Info("Handling %s\n", path)

	switch path {
	case "/storage":
		a.StorageHandler.ServeHTTP(w, r)
	case "/ping":
		a.PingHandler.ServeHTTP(w, r)
	}
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}