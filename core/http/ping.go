package http

import "net/http"

type PingService struct {

}

func (p *PingService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
