package dto

type StatsResponse struct {
	Jobs    JobStats    `json:"jobs"`
	Workers WorkerStats `json:"workers"`
}

type JobStats struct {
	Pending   int `json:"pending"`
	Running   int `json:"running"`
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
}

type WorkerStats struct {
	Total int `json:"total"`
	Idle  int `json:"idle"`
	Busy  int `json:"busy"`
}
