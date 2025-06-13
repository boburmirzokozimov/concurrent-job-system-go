package worker

import (
	"concurrent-job-system/models"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

type JobLogger struct {
	info  *log.Logger
	fail  *log.Logger
	start *log.Logger
	retry *log.Logger
}

var jobLog = &JobLogger{
	info:  log.New(os.Stdout, fmt.Sprintf("%s[INFO]%s ", ColorGreen, ColorReset), log.Ltime),
	fail:  log.New(os.Stdout, fmt.Sprintf("%s[FAIL]%s ", ColorRed, ColorReset), log.Ltime),
	start: log.New(os.Stdout, fmt.Sprintf("%s[START]%s ", ColorCyan, ColorReset), log.Ltime),
	retry: log.New(os.Stdout, fmt.Sprintf("%s[RETRY]%s ", ColorYellow, ColorReset), log.Ltime),
}

func LogJobStart(job models.Processable) {
	priorityColor := models.ColorForPriority(job.GetPriority().String())
	jobLog.start.Printf("%s[%s][%s] Job ID: %d | Retry: %d/%d%s",
		priorityColor, job.Type(), job.GetPriority(), job.GetId(), job.GetRetries(), job.GetMaxRetryCount(), ColorReset)
}

func LogJobSuccess(job models.Processable) {
	priorityColor := models.ColorForPriority(job.GetPriority().String())
	jobLog.info.Printf("%s[%s][%s] Job ID: %d completed successfully%s",
		priorityColor, job.Type(), job.GetPriority(), job.GetId(), ColorReset)
}

func LogJobRetry(job models.Processable, backoff time.Duration) {
	priorityColor := models.ColorForPriority(job.GetPriority().String())
	jobLog.retry.Printf("%s[%s][%s] Job ID: %d retrying in %v...%s",
		priorityColor, job.Type(), job.GetPriority(), job.GetId(), backoff, ColorReset)
}

func LogJobFail(job models.Processable) {
	priorityColor := models.ColorForPriority(job.GetPriority().String())
	jobLog.fail.Printf("%s[%s][%s] Job ID: %d failed%s",
		priorityColor, job.Type(), job.GetPriority(), job.GetId(), ColorReset)
}

func LogJobCanceled(job models.Processable) {
	priorityColor := models.ColorForPriority(job.GetPriority().String())
	jobLog.fail.Printf("%s[%s][%s] Job ID: %d canceled by context%s",
		priorityColor, job.Type(), job.GetPriority(), job.GetId(), ColorReset)
}
