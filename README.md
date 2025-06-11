# Concurrent Job Processing System in Go 🧵⚙️

This project is a concurrency-focused job processing system built in **Go**, demonstrating advanced use of:

- Goroutines
- Worker pools
- Context cancellation
- Exponential backoff retry
- Sync primitives (`sync.WaitGroup`, `sync.Map`, `atomic`)

---

## 🚀 Features

- Bounded worker pool
- Graceful shutdown via OS signals (`SIGINT`, `SIGTERM`)
- Context-aware cancelable jobs
- Exponential backoff retry logic
- Real-time stats tracking with atomic counters
- Thread-safe per-job status map