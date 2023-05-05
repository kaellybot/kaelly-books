package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/repositories/alignments"
)

func New(broker amqp.MessageBroker, alignBookRepo alignments.Repository) *Impl {
	return &Impl{
		broker:        broker,
		alignBookRepo: alignBookRepo,
	}
}
