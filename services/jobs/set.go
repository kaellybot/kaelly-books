package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/kaellybot/kaelly-books/utils/replies"
	"github.com/rs/zerolog/log"
)

func (service *Impl) SetRequest(ctx amqp.Context, request *amqp.JobSetRequest,
	game amqp.Game, lg amqp.Language) {
	if !isValidJobSetRequest(request) {
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_JOB_SET_ANSWER, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogUserID, request.UserId).
		Str(constants.LogJobID, request.JobId).
		Str(constants.LogServerID, request.ServerId).
		Msgf("Set job request received")

	jobBook := entities.JobBook{
		UserID:   request.UserId,
		JobID:    request.JobId,
		ServerID: request.ServerId,
		Game:     game,
		Level:    request.Level,
	}

	var err error
	if request.Level > 0 {
		err = service.jobBookRepo.SaveUserBook(jobBook)
	} else {
		err = service.jobBookRepo.DeleteUserBook(jobBook)
	}
	if err != nil {
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_JOB_SET_ANSWER, lg)
		return
	}

	answer := mappers.MapJobSetAnswer(lg)
	replies.SucceededAnswer(ctx, service.broker, answer)
}

func isValidJobSetRequest(request *amqp.JobSetRequest) bool {
	return request != nil
}
