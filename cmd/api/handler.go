package api

import (
	"concurrent-job-system/internal/container"
	"concurrent-job-system/internal/job"
	"concurrent-job-system/internal/job/factory"
	"concurrent-job-system/pkg/dto"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	c *container.Container
}

func NewHandler(c *container.Container) *Handler {
	return &Handler{c}
}

func (h *Handler) SubmitJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SubmitJobRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	j := h.MakeJobFromPayload(req)
	if j == nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	h.c.Pool.Submit(j)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) // Optional: use 202 Accepted
	_ = json.NewEncoder(w).Encode(map[string]any{
		"id":     j.GetId(),
		"status": "queued",
	})
}

func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/jobs/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	status, err := h.c.JobRepository.Load(id)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	resp := JobResponse{ID: id, Status: status.GetStatus()}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.c.JobRepository.GetStats()
	if err != nil {
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	workerStats := h.c.Pool.GetStats()

	resp := dto.StatsResponse{
		Jobs:    stats,
		Workers: workerStats,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) MakeJobFromPayload(r SubmitJobRequest) job.IProcessable {
	payloadJson, err := json.Marshal(r.Payload)
	if err != nil {
		h.c.Logger.Error("Invalid payload: %v", err)
		return nil
	}

	row := job.BaseJob{
		JobType: r.Type,
		Payload: string(payloadJson),
	}

	j, err := factory.DeserializeJob(row)
	if err != nil {
		h.c.Logger.Warn("Failed to deserialize job: %v", err)
		return nil
	}
	return j
}
