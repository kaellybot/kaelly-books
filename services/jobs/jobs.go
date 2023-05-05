package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/repositories/jobs"
)

func New(broker amqp.MessageBroker, jobBookRepo jobs.Repository) *Impl {
	return &Impl{
		broker:      broker,
		jobBookRepo: jobBookRepo,
	}
}
