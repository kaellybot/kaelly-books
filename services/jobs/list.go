package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/kaellybot/kaelly-books/utils/replies"
	"github.com/rs/zerolog/log"
)

func (service *Impl) GetBookRequest(ctx amqp.Context, message *amqp.RabbitMQMessage) {
	request := message.GetJobGetBookRequest()
	if !isValidJobGetRequest(request) {
		replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_JOB_GET_BOOK_ANSWER)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogGame, message.GetGame().String()).
		Str(constants.LogJobID, request.GetJobId()).
		Str(constants.LogServerID, request.GetServerId()).
		Msgf("Get job books request received")

	books, total, err := service.jobBookRepo.GetBooks(request.GetJobId(), request.GetServerId(),
		request.GetUserIds(), message.GetGame(), int(request.GetOffset()), int(request.GetSize()))
	if err != nil {
		replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_JOB_GET_BOOK_ANSWER)
		return
	}

	answer := mappers.MapJobBookAnswer(message, books, total)
	replies.SucceededAnswer(ctx, service.broker, answer)
}

func isValidJobGetRequest(request *amqp.JobGetBookRequest) bool {
	return request != nil
}
