package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/kaellybot/kaelly-books/utils/replies"
	"github.com/rs/zerolog/log"
)

func (service *Impl) SetRequest(ctx amqp.Context, request *amqp.AlignSetRequest,
	lg amqp.Language) {
	if !isValidAlignSetRequest(request) {
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_ALIGN_SET_ANSWER, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogUserID, request.UserId).
		Str(constants.LogCityID, request.CityId).
		Str(constants.LogOrderID, request.OrderId).
		Str(constants.LogServerID, request.ServerId).
		Msgf("Set job request received")

	jobBook := entities.AlignmentBook{
		UserID:   request.UserId,
		CityID:   request.CityId,
		OrderID:  request.OrderId,
		ServerID: request.ServerId,
		Level:    request.Level,
	}

	var err error
	if request.Level > 0 {
		err = service.alignBookRepo.SaveUserBook(jobBook)
	} else {
		err = service.alignBookRepo.DeleteUserBook(jobBook)
	}
	if err != nil {
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_ALIGN_SET_ANSWER, lg)
		return
	}

	answer := mappers.MapAlignSetAnswer(lg)
	replies.SucceededAnswer(ctx, service.broker, answer)
}

func isValidAlignSetRequest(request *amqp.AlignSetRequest) bool {
	return request != nil
}
