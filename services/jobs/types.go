package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/repositories/jobs"
)

type Service interface {
	GetBookRequest(request *amqp.JobGetBookRequest, correlationID, answersRoutingkey string, lg amqp.Language)
	SetRequest(request *amqp.JobSetRequest, correlationID, answersRoutingkey string, lg amqp.Language)
	UserRequest(request *amqp.JobGetUserRequest, correlationID, answersRoutingkey string, lg amqp.Language)
}

type Impl struct {
	broker      amqp.MessageBroker
	jobBookRepo jobs.Repository
}
