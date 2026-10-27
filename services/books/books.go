package books

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/services/alignments"
	"github.com/kaellybot/kaelly-books/services/jobs"
	"github.com/kaellybot/kaelly-books/utils/replies"
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
	switch message.GetType() {
	case amqp.RabbitMQMessage_JOB_GET_BOOK_REQUEST:
		service.handle(ctx, message, amqp.RabbitMQMessage_JOB_GET_BOOK_ANSWER,
			service.jobService.GetBookRequest)
	case amqp.RabbitMQMessage_JOB_GET_USER_REQUEST:
		service.handle(ctx, message, amqp.RabbitMQMessage_JOB_GET_USER_ANSWER,
			service.jobService.UserRequest)
	case amqp.RabbitMQMessage_JOB_SET_REQUEST:
		service.handle(ctx, message, amqp.RabbitMQMessage_JOB_SET_ANSWER,
			service.jobService.SetRequest)
	case amqp.RabbitMQMessage_ALIGN_GET_BOOK_REQUEST:
		service.handle(ctx, message, amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER,
			service.alignService.GetBookRequest)
	case amqp.RabbitMQMessage_ALIGN_GET_USER_REQUEST:
		service.handle(ctx, message, amqp.RabbitMQMessage_ALIGN_GET_USER_ANSWER,
			service.alignService.UserRequest)
	case amqp.RabbitMQMessage_ALIGN_SET_REQUEST:
		service.handle(ctx, message, amqp.RabbitMQMessage_ALIGN_SET_ANSWER,
			service.alignService.SetRequest)
	default:
		log.Warn().
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Msgf("Type not recognized, request ignored")
	}
}

// handle validates the request game before running the handler, so that a request
// never gets another game's data.
func (service *Impl) handle(ctx amqp.Context, message *amqp.RabbitMQMessage,
	answerType amqp.RabbitMQMessage_Type, handler requestHandler) {
	// The broker guards already refuse these, so reaching this point means the guard
	// was bypassed. Answering is impossible: the reply would carry ANY_GAME too.
	if !message.IsGameSet() {
		log.Error().
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Msgf("Request without game received, request ignored")
		return
	}

	if !isGameSupported(message.GetGame()) {
		log.Warn().
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Str(constants.LogGame, message.GetGame().String()).
			Msgf("Game not supported, request refused")
		replies.FailedAnswer(ctx, service.broker, message, answerType)
		return
	}

	handler(ctx, message)
}

// isGameSupported reports whether books can serve this game. Books are game-generic:
// a game is supported as soon as its jobs, cities and orders exist in database.
func isGameSupported(game amqp.Game) bool {
	switch game {
	case amqp.Game_DOFUS_GAME, amqp.Game_DOFUS_TOUCH, amqp.Game_DOFUS_RETRO:
		return true
	case amqp.Game_ANY_GAME:
		return false
	default:
		return false
	}
}
