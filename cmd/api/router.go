package api

import "net/http"

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/jobs", h.SubmitJob)
	mux.HandleFunc("/jobs/", h.GetJobStatus)
	return mux
}
