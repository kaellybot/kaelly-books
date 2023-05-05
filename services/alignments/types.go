package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/repositories/alignments"
)

type Service interface {
	GetBookRequest(request *amqp.AlignGetBookRequest, correlationID, answersRoutingkey string, lg amqp.Language)
	SetRequest(request *amqp.AlignSetRequest, correlationID, answersRoutingkey string, lg amqp.Language)
	UserRequest(request *amqp.AlignGetUserRequest, correlationID, answersRoutingkey string, lg amqp.Language)
}

type Impl struct {
	broker        amqp.MessageBroker
	alignBookRepo alignments.Repository
}
