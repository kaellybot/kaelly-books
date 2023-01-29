package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/repositories/jobs"
)

func New(broker amqp.MessageBrokerInterface, jobBookRepo jobs.JobBookRepository) *JobServiceImpl {
	return &JobServiceImpl{
		broker:      broker,
		jobBookRepo: jobBookRepo,
	}
}
