package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/repositories/alignments"
)

type Service interface {
	GetBookRequest(ctx amqp.Context, message *amqp.RabbitMQMessage)
	SetRequest(ctx amqp.Context, message *amqp.RabbitMQMessage)
	UserRequest(ctx amqp.Context, message *amqp.RabbitMQMessage)
}

type Impl struct {
	broker        amqp.MessageBroker
	alignBookRepo alignments.Repository
}
