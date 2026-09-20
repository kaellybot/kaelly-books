package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/kaellybot/kaelly-books/utils/replies"
	"github.com/rs/zerolog/log"
)

func (service *Impl) SetRequest(ctx amqp.Context, message *amqp.RabbitMQMessage) {
	request := message.GetJobSetRequest()
	if !isValidJobSetRequest(request) {
		replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_JOB_SET_ANSWER)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogGame, message.GetGame().String()).
		Str(constants.LogUserID, request.GetUserId()).
		Str(constants.LogJobID, request.GetJobId()).
		Str(constants.LogServerID, request.GetServerId()).
		Msgf("Set job request received")

	jobBook := entities.JobBook{
		UserID:   request.GetUserId(),
		JobID:    request.GetJobId(),
		ServerID: request.GetServerId(),
		Game:     message.GetGame(),
		Level:    request.GetLevel(),
	}

	var err error
	if request.GetLevel() > 0 {
		err = service.jobBookRepo.SaveUserBook(jobBook)
	} else {
		err = service.jobBookRepo.DeleteUserBook(jobBook)
	}
	if err != nil {
		replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_JOB_SET_ANSWER)
		return
	}

	answer := mappers.MapJobSetAnswer(message)
	replies.SucceededAnswer(ctx, service.broker, answer)
}

func isValidJobSetRequest(request *amqp.JobSetRequest) bool {
	return request != nil
}
