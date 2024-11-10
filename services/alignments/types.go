package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/repositories/alignments"
)

type Service interface {
	GetBookRequest(ctx amqp.Context, request *amqp.AlignGetBookRequest, lg amqp.Language)
	SetRequest(ctx amqp.Context, request *amqp.AlignSetRequest, lg amqp.Language)
	UserRequest(ctx amqp.Context, request *amqp.AlignGetUserRequest, lg amqp.Language)
}

type Impl struct {
	broker        amqp.MessageBroker
	alignBookRepo alignments.Repository
}
