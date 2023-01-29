package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/repositories/alignments"
)

func New(broker amqp.MessageBrokerInterface, alignBookRepo alignments.AlignmentBookRepository) *AlignmentServiceImpl {
	return &AlignmentServiceImpl{
		broker:        broker,
		alignBookRepo: alignBookRepo,
	}
}
