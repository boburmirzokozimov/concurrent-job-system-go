package worker

import (
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
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[97m"

	// ColorBrightRed Vivid background-like text highlights for priorities
	ColorBrightRed    = "\033[1;91m"
	ColorBrightGreen  = "\033[1;92m"
	ColorBrightYellow = "\033[1;93m"
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

func LogJobStart(job Processable) {
	priorityColor := ColorForPriority(job.GetPriority())
	jobLog.start.Printf("%s[%s][%s] Job ID: %d | Retry: %d/%d%s",
		priorityColor, job.Type(), job.GetPriority().String(), job.GetId(), job.GetRetries(), job.GetMaxRetryCount(), ColorReset)
}

func LogJobSuccess(job Processable) {
	priorityColor := ColorForPriority(job.GetPriority())
	jobLog.info.Printf("%s[%s][%s] Job ID: %d completed successfully%s",
		priorityColor, job.Type(), job.GetPriority().String(), job.GetId(), ColorReset)
}

func LogJobRetry(job Processable, backoff time.Duration) {
	priorityColor := ColorForPriority(job.GetPriority())
	jobLog.retry.Printf("%s[%s][%s] Job ID: %d retrying in %v...%s",
		priorityColor, job.Type(), job.GetPriority().String(), job.GetId(), backoff, ColorReset)
}

func LogJobFail(job Processable) {
	priorityColor := ColorForPriority(job.GetPriority())
	jobLog.fail.Printf("%s[%s][%s] Job ID: %d failed%s",
		priorityColor, job.Type(), job.GetPriority().String(), job.GetId(), ColorReset)
}

func LogJobCanceled(job Processable) {
	priorityColor := ColorForPriority(job.GetPriority())
	jobLog.fail.Printf("%s[%s][%s] Job ID: %d canceled by context%s",
		priorityColor, job.Type(), job.GetPriority().String(), job.GetId(), ColorReset)
}
