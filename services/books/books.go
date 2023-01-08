package books

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/repositories/jobs"
	"github.com/rs/zerolog/log"
)

func New(broker amqp.MessageBrokerInterface, jobBookRepo jobs.JobBookRepository) (*BooksServiceImpl, error) {

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

func (service *BooksServiceImpl) consume(ctx context.Context, message *amqp.RabbitMQMessage, correlationId string) {
	// TODO
}
