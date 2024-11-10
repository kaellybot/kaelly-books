package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/kaellybot/kaelly-books/utils/replies"
	"github.com/rs/zerolog/log"
)

func (service *Impl) GetBookRequest(ctx amqp.Context, request *amqp.JobGetBookRequest,
	lg amqp.Language) {
	if !isValidJobGetRequest(request) {
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_JOB_GET_BOOK_ANSWER, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogJobID, request.JobId).
		Str(constants.LogServerID, request.ServerId).
		Msgf("Get job books request received")

	books, total, err := service.jobBookRepo.GetBooks(request.GetJobId(), request.GetServerId(),
		request.GetUserIds(), int(request.GetOffset()), int(request.GetSize()))
	if err != nil {
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_JOB_GET_BOOK_ANSWER, lg)
		return
	}

	answer := mappers.MapJobBookAnswer(request, books, total, lg)
	replies.SucceededAnswer(ctx, service.broker, answer)
}

func isValidJobGetRequest(request *amqp.JobGetBookRequest) bool {
	return request != nil
}
