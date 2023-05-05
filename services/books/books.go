package books

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/services/alignments"
	"github.com/kaellybot/kaelly-books/services/jobs"
	"github.com/rs/zerolog/log"
)

func New(broker amqp.MessageBroker, jobService jobs.Service,
	alignService alignments.Service) *Impl {
	return &Impl{
		broker:       broker,
		jobService:   jobService,
		alignService: alignService,
	}
}

func GetBinding() amqp.Binding {
	return amqp.Binding{
		Exchange:   amqp.ExchangeRequest,
		RoutingKey: requestsRoutingkey,
		Queue:      requestQueueName,
	}
}

func (service *Impl) Consume() error {
	log.Info().Msgf("Consuming books requests...")
	return service.broker.Consume(requestQueueName, service.consume)
}

func (service *Impl) consume(_ context.Context,
	message *amqp.RabbitMQMessage, correlationID string) {
	//exhaustive:ignore Don't need to be exhaustive here since they will be handled by default case
	switch message.Type {
	case amqp.RabbitMQMessage_JOB_GET_BOOK_REQUEST:
		service.jobService.GetBookRequest(message.JobGetBookRequest, correlationID, answersRoutingkey, message.Language)
	case amqp.RabbitMQMessage_JOB_GET_USER_REQUEST:
		service.jobService.UserRequest(message.JobGetUserRequest, correlationID, answersRoutingkey, message.Language)
	case amqp.RabbitMQMessage_JOB_SET_REQUEST:
		service.jobService.SetRequest(message.JobSetRequest, correlationID, answersRoutingkey, message.Language)
	case amqp.RabbitMQMessage_ALIGN_GET_BOOK_REQUEST:
		service.alignService.GetBookRequest(message.AlignGetBookRequest, correlationID, answersRoutingkey, message.Language)
	case amqp.RabbitMQMessage_ALIGN_GET_USER_REQUEST:
		service.alignService.UserRequest(message.AlignGetUserRequest, correlationID, answersRoutingkey, message.Language)
	case amqp.RabbitMQMessage_ALIGN_SET_REQUEST:
		service.alignService.SetRequest(message.AlignSetRequest, correlationID, answersRoutingkey, message.Language)
	default:
		log.Warn().
			Str(constants.LogCorrelationID, correlationID).
			Msgf("Type not recognized, request ignored")
	}
}
