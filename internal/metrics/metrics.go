package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	QueuedJobs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "job_queued_total",
			Help: "Total jobs queued by type",
		},
		[]string{"type"},
	)

	RunningJobs = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "job_running",
			Help: "Number of jobs currently running",
		})

	SucceededJobs = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "job_success_total",
			Help: "Total number of successfully completed jobs",
		})

	FailedJobs = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "job_failed_total",
			Help: "Total number of failed jobs",
		})

	DeadLetterJobs = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "job_deadletter_total",
			Help: "Total number of jobs sent to dead-letter",
		})

	JobLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "job_latency_seconds",
			Help:    "Job execution duration in seconds",
			Buckets: prometheus.DefBuckets,
		})

	QueuedGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "job_queue_length",
			Help: "Current number of jobs waiting in the queue",
		})
)

func Init() {
	prometheus.MustRegister(
		QueuedJobs,
		RunningJobs,
		SucceededJobs,
		FailedJobs,
		DeadLetterJobs,
		JobLatency,
		QueuedGauge,
	)
}
