package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/kaellybot/kaelly-books/utils/replies"
	"github.com/rs/zerolog/log"
)

func (service *Impl) SetRequest(ctx amqp.Context, message *amqp.RabbitMQMessage) {
	request := message.GetAlignSetRequest()
	if !isValidAlignSetRequest(request) {
		replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_ALIGN_SET_ANSWER)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogGame, message.GetGame().String()).
		Str(constants.LogUserID, request.GetUserId()).
		Str(constants.LogCityID, request.GetCityId()).
		Str(constants.LogOrderID, request.GetOrderId()).
		Str(constants.LogServerID, request.GetServerId()).
		Msgf("Set align request received")

	alignBook := entities.AlignmentBook{
		UserID:   request.GetUserId(),
		CityID:   request.GetCityId(),
		OrderID:  request.GetOrderId(),
		ServerID: request.GetServerId(),
		Game:     message.GetGame(),
		Level:    request.GetLevel(),
	}

	var err error
	if request.GetLevel() > 0 {
		err = service.alignBookRepo.SaveUserBook(alignBook)
	} else {
		err = service.alignBookRepo.DeleteUserBook(alignBook)
	}
	if err != nil {
		replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_ALIGN_SET_ANSWER)
		return
	}

	answer := mappers.MapAlignSetAnswer(message)
	replies.SucceededAnswer(ctx, service.broker, answer)
}

func isValidAlignSetRequest(request *amqp.AlignSetRequest) bool {
	return request != nil
}
