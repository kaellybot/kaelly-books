package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/repositories/alignments"
)

type AlignmentService interface {
	GetBookRequest(request *amqp.AlignGetBookRequest, correlationId, answersRoutingkey string, lg amqp.Language)
	SetRequest(request *amqp.AlignSetRequest, correlationId, answersRoutingkey string, lg amqp.Language)
	UserRequest(request *amqp.AlignGetUserRequest, correlationId, answersRoutingkey string, lg amqp.Language)
}

type AlignmentServiceImpl struct {
	broker        amqp.MessageBrokerInterface
	alignBookRepo alignments.AlignmentBookRepository
}
