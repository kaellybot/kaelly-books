package books

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/models/constants"
	"github.com/kaellybot/kaelly-configurator/repositories/jobs"
	"github.com/rs/zerolog/log"
)

func New(broker amqp.MessageBrokerInterface, jobBookRepo jobs.JobBookRepository) (
	*BooksServiceImpl, error) {

	return &BooksServiceImpl{
		broker:      broker,
		jobBookRepo: jobBookRepo,
	}, nil
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
		service.getBookRequest(message, correlationId)
	case amqp.RabbitMQMessage_JOB_GET_USER_REQUEST:
		service.userRequest(message, correlationId)
	case amqp.RabbitMQMessage_JOB_SET_REQUEST:
		service.setRequest(message, correlationId)
	default:
		log.Warn().
			Str(constants.LogCorrelationId, correlationId).
			Msgf("Type not recognized, request ignored")
	}
}
