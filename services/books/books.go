package books

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/models/constants"
	"github.com/kaellybot/kaelly-configurator/services/alignments"
	"github.com/kaellybot/kaelly-configurator/services/jobs"
	"github.com/rs/zerolog/log"
)

func New(broker amqp.MessageBrokerInterface, jobService jobs.JobService,
	alignService alignments.AlignmentService) *BooksServiceImpl {

	return &BooksServiceImpl{
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

func (service *BooksServiceImpl) Consume() error {
	log.Info().Msgf("Consuming books requests...")
	return service.broker.Consume(requestQueueName, requestsRoutingkey, service.consume)
}

func (service *BooksServiceImpl) consume(ctx context.Context,
	message *amqp.RabbitMQMessage, correlationId string) {

	switch message.Type {
	case amqp.RabbitMQMessage_JOB_GET_BOOK_REQUEST:
		service.jobService.GetBookRequest(message.JobGetBookRequest, correlationId, answersRoutingkey, message.Language)
	case amqp.RabbitMQMessage_JOB_GET_USER_REQUEST:
		service.jobService.UserRequest(message.JobGetUserRequest, correlationId, answersRoutingkey, message.Language)
	case amqp.RabbitMQMessage_JOB_SET_REQUEST:
		service.jobService.SetRequest(message.JobSetRequest, correlationId, answersRoutingkey, message.Language)
	case amqp.RabbitMQMessage_ALIGN_GET_BOOK_REQUEST:
		service.alignService.GetBookRequest(message.AlignGetBookRequest, correlationId, answersRoutingkey, message.Language)
	case amqp.RabbitMQMessage_ALIGN_GET_USER_REQUEST:
		service.alignService.UserRequest(message.AlignGetUserRequest, correlationId, answersRoutingkey, message.Language)
	case amqp.RabbitMQMessage_ALIGN_SET_REQUEST:
		service.alignService.SetRequest(message.AlignSetRequest, correlationId, answersRoutingkey, message.Language)
	default:
		log.Warn().
			Str(constants.LogCorrelationId, correlationId).
			Msgf("Type not recognized, request ignored")
	}
}
