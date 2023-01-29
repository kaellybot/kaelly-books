package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/repositories/jobs"
)

type JobService interface {
	GetBookRequest(request *amqp.JobGetBookRequest, correlationId, answersRoutingkey string, lg amqp.Language)
	SetRequest(request *amqp.JobSetRequest, correlationId, answersRoutingkey string, lg amqp.Language)
	UserRequest(request *amqp.JobGetUserRequest, correlationId, answersRoutingkey string, lg amqp.Language)
}

type JobServiceImpl struct {
	broker      amqp.MessageBrokerInterface
	jobBookRepo jobs.JobBookRepository
}
