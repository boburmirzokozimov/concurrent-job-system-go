package models

type Processable interface {
	GetRetries() int
	GetId() int
	GetMaxRetryCount() int
	Process() error
	IncRetry()
	Type() string
	GetPriority() Priority
	Base() *BaseJob
	PayloadOnly() interface{}
}

type JobPayload interface {
	ToProcessable(base BaseJob) Processable
}
