package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/kaellybot/kaelly-books/utils/replies"
	"github.com/rs/zerolog/log"
)

func (service *Impl) UserRequest(ctx amqp.Context, request *amqp.AlignGetUserRequest,
	game amqp.Game, lg amqp.Language) {
	if !isValidAlignGetUserRequest(request) {
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_ALIGN_GET_USER_ANSWER, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogUserID, request.UserId).
		Str(constants.LogServerID, request.ServerId).
		Msgf("Get job user request received")

	books, err := service.alignBookRepo.GetUserBook(request.UserId, request.ServerId, game)
	if err != nil {
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_ALIGN_GET_USER_ANSWER, lg)
		return
	}

	answer := mappers.MapAlignUserAnswer(books, request.ServerId, lg)
	replies.SucceededAnswer(ctx, service.broker, answer)
}

func isValidAlignGetUserRequest(request *amqp.AlignGetUserRequest) bool {
	return request != nil
}
