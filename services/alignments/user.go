package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/kaellybot/kaelly-books/utils/replies"
	"github.com/rs/zerolog/log"
)

func (service *Impl) UserRequest(ctx amqp.Context, message *amqp.RabbitMQMessage) {
	request := message.GetAlignGetUserRequest()
	if !isValidAlignGetUserRequest(request) {
		replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_ALIGN_GET_USER_ANSWER)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogGame, message.GetGame().String()).
		Str(constants.LogUserID, request.GetUserId()).
		Str(constants.LogServerID, request.GetServerId()).
		Msgf("Get align user request received")

	books, err := service.alignBookRepo.GetUserBook(request.GetUserId(), request.GetServerId(),
		message.GetGame())
	if err != nil {
		replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_ALIGN_GET_USER_ANSWER)
		return
	}

	answer := mappers.MapAlignUserAnswer(message, books)
	replies.SucceededAnswer(ctx, service.broker, answer)
}

func isValidAlignGetUserRequest(request *amqp.AlignGetUserRequest) bool {
	return request != nil
}
