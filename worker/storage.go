package worker

import (
	"encoding/json"
	"os"
	"sync"
)

type FileJobStorage struct {
	path string
	mu   sync.Mutex
}

func NewFileJobStorage(path string) *FileJobStorage {
	return &FileJobStorage{
		path: path,
	}
}
func (s *FileJobStorage) Save(j Processable) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	jobs := s.LoadAll()
	jobs = append(jobs, *j.Base())

	data, _ := json.MarshalIndent(jobs, "", "  ")
	return os.WriteFile(s.path, data, os.ModePerm)
}

func (s *FileJobStorage) LoadAll() []BaseJob {

	file, err := os.ReadFile(s.path)
	if err != nil {
		return []BaseJob{}
	}
	var jobs []BaseJob
	_ = json.Unmarshal(file, &jobs)
	return jobs
}

func (s *FileJobStorage) MarkCompleted(id int, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	jobs := s.LoadAll()
	for i := range jobs {
		if jobs[i].ID == id {
			jobs[i].Status = status
		}
	}
	data, _ := json.MarshalIndent(jobs, "", "  ")
	return os.WriteFile(s.path, data, os.ModePerm)
}

func (s *FileJobStorage) LoadPending() []BaseJob {
	s.mu.Lock()
	defer s.mu.Unlock()
	var pending []BaseJob
	jobs := s.LoadAll()
	for i := range jobs {
		if jobs[i].Status != "success" {
			pending = append(pending, jobs[i])
		}
	}

	return pending
}
