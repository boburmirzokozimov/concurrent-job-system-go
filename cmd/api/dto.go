package api

type SubmitJobRequest struct {
	Type     string `json:"type"`
	Priority int    `json:"priority"`
	Payload  any    `json:"payload"`
}

type JobResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}
