package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/kaellybot/kaelly-books/utils/replies"
	"github.com/rs/zerolog/log"
)

func (service *Impl) GetBookRequest(ctx amqp.Context, request *amqp.AlignGetBookRequest,
	lg amqp.Language) {
	if !isValidAlignGetRequest(request) {
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogCityID, request.CityId).
		Str(constants.LogOrderID, request.OrderId).
		Str(constants.LogServerID, request.ServerId).
		Msgf("Get align books request received")

	books, total, err := service.alignBookRepo.GetBooks(request.GetCityId(), request.GetOrderId(),
		request.GetServerId(), request.GetUserIds(), int(request.GetOffset()), int(request.GetSize()))
	if err != nil {
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER, lg)
		return
	}

	answer := mappers.MapAlignBookAnswer(request, books, total, lg)
	replies.SucceededAnswer(ctx, service.broker, answer)
}

func isValidAlignGetRequest(request *amqp.AlignGetBookRequest) bool {
	return request != nil
}
