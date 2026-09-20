package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/kaellybot/kaelly-books/utils/replies"
	"github.com/rs/zerolog/log"
)

func (service *Impl) GetBookRequest(ctx amqp.Context, message *amqp.RabbitMQMessage) {
	request := message.GetAlignGetBookRequest()
	if !isValidAlignGetRequest(request) {
		replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogGame, message.GetGame().String()).
		Str(constants.LogCityID, request.GetCityId()).
		Str(constants.LogOrderID, request.GetOrderId()).
		Str(constants.LogServerID, request.GetServerId()).
		Msgf("Get align books request received")

	books, total, err := service.alignBookRepo.GetBooks(request.GetCityId(), request.GetOrderId(),
		request.GetServerId(), request.GetUserIds(), message.GetGame(),
		int(request.GetOffset()), int(request.GetSize()))
	if err != nil {
		replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER)
		return
	}

	answer := mappers.MapAlignBookAnswer(message, books, total)
	replies.SucceededAnswer(ctx, service.broker, answer)
}

func isValidAlignGetRequest(request *amqp.AlignGetBookRequest) bool {
	return request != nil
}
