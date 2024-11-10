package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/kaellybot/kaelly-books/utils/replies"
	"github.com/rs/zerolog/log"
)

func (service *Impl) UserRequest(ctx amqp.Context, request *amqp.JobGetUserRequest,
	lg amqp.Language) {
	if !isValidJobGetUserRequest(request) {
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_JOB_GET_USER_ANSWER, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogUserID, request.UserId).
		Str(constants.LogServerID, request.ServerId).
		Msgf("Get job user request received")

	books, err := service.jobBookRepo.GetUserBook(request.UserId, request.ServerId)
	if err != nil {
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_JOB_GET_USER_ANSWER, lg)
		return
	}

	answer := mappers.MapJobUserAnswer(books, request.ServerId, lg)
	replies.SucceededAnswer(ctx, service.broker, answer)
}

func isValidJobGetUserRequest(request *amqp.JobGetUserRequest) bool {
	return request != nil
}
