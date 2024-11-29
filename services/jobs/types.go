package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/repositories/jobs"
)

type Service interface {
	GetBookRequest(ctx amqp.Context, request *amqp.JobGetBookRequest, game amqp.Game, lg amqp.Language)
	SetRequest(ctx amqp.Context, request *amqp.JobSetRequest, game amqp.Game, lg amqp.Language)
	UserRequest(ctx amqp.Context, request *amqp.JobGetUserRequest, game amqp.Game, lg amqp.Language)
}

type Impl struct {
	broker      amqp.MessageBroker
	jobBookRepo jobs.Repository
}
