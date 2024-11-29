package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/repositories/alignments"
)

type Service interface {
	GetBookRequest(ctx amqp.Context, request *amqp.AlignGetBookRequest, game amqp.Game, lg amqp.Language)
	SetRequest(ctx amqp.Context, request *amqp.AlignSetRequest, game amqp.Game, lg amqp.Language)
	UserRequest(ctx amqp.Context, request *amqp.AlignGetUserRequest, game amqp.Game, lg amqp.Language)
}

type Impl struct {
	broker        amqp.MessageBroker
	alignBookRepo alignments.Repository
}
