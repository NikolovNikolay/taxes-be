package echoaddons

import "net/http"

type MaxBytesHandler struct {
	Handler      http.Handler
	MaxReqLength int64
}

func (h *MaxBytesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.MaxReqLength)
	h.Handler.ServeHTTP(w, r)
}
