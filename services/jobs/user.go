package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/kaellybot/kaelly-books/utils/replies"
	"github.com/rs/zerolog/log"
)

func (service *Impl) UserRequest(ctx amqp.Context, message *amqp.RabbitMQMessage) {
	request := message.GetJobGetUserRequest()
	if !isValidJobGetUserRequest(request) {
		replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_JOB_GET_USER_ANSWER)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogGame, message.GetGame().String()).
		Str(constants.LogUserID, request.GetUserId()).
		Str(constants.LogServerID, request.GetServerId()).
		Msgf("Get job user request received")

	books, err := service.jobBookRepo.GetUserBook(request.GetUserId(), request.GetServerId(),
		message.GetGame())
	if err != nil {
		replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_JOB_GET_USER_ANSWER)
		return
	}

	answer := mappers.MapJobUserAnswer(message, books)
	replies.SucceededAnswer(ctx, service.broker, answer)
}

func isValidJobGetUserRequest(request *amqp.JobGetUserRequest) bool {
	return request != nil
}
