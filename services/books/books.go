package books

import (
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

func (service *Impl) Consume() {
	log.Info().Msgf("Consuming books requests...")
	service.broker.Consume(requestQueueName, service.consume)
}

func (service *Impl) consume(ctx amqp.Context, message *amqp.RabbitMQMessage) {
	//exhaustive:ignore Don't need to be exhaustive here since they will be handled by default case
	switch message.Type {
	case amqp.RabbitMQMessage_JOB_GET_BOOK_REQUEST:
		service.jobService.GetBookRequest(ctx, message.JobGetBookRequest, message.Game, message.Language)
	case amqp.RabbitMQMessage_JOB_GET_USER_REQUEST:
		service.jobService.UserRequest(ctx, message.JobGetUserRequest, message.Game, message.Language)
	case amqp.RabbitMQMessage_JOB_SET_REQUEST:
		service.jobService.SetRequest(ctx, message.JobSetRequest, message.Game, message.Language)
	case amqp.RabbitMQMessage_ALIGN_GET_BOOK_REQUEST:
		service.alignService.GetBookRequest(ctx, message.AlignGetBookRequest, message.Game, message.Language)
	case amqp.RabbitMQMessage_ALIGN_GET_USER_REQUEST:
		service.alignService.UserRequest(ctx, message.AlignGetUserRequest, message.Game, message.Language)
	case amqp.RabbitMQMessage_ALIGN_SET_REQUEST:
		service.alignService.SetRequest(ctx, message.AlignSetRequest, message.Game, message.Language)
	default:
		log.Warn().
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Msgf("Type not recognized, request ignored")
	}
}
