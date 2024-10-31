package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/rs/zerolog/log"
)

func (service *Impl) GetBookRequest(request *amqp.AlignGetBookRequest,
	correlationID, answersRoutingkey string, lg amqp.Language) {
	if !isValidAlignGetRequest(request) {
		service.publishFailedGetBookAnswer(correlationID, answersRoutingkey, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Str(constants.LogCityID, request.CityId).
		Str(constants.LogOrderID, request.OrderId).
		Str(constants.LogServerID, request.ServerId).
		Msgf("Get job books request received")

	books, total, err := service.alignBookRepo.GetBooks(request.GetCityId(), request.GetOrderId(),
		request.GetServerId(), request.GetUserIds(), int(request.GetOffset()), int(request.GetSize()))
	if err != nil {
		service.publishFailedGetBookAnswer(correlationID, answersRoutingkey, lg)
		return
	}

	answer := mappers.MapAlignBookAnswer(request, books, total)
	service.publishSucceededGetBookAnswer(correlationID, answersRoutingkey, answer, lg)
}

func (service *Impl) publishSucceededGetBookAnswer(correlationID, answersRoutingkey string,
	answer *amqp.AlignGetBookAnswer, lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:               amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER,
		Status:             amqp.RabbitMQMessage_SUCCESS,
		Language:           lg,
		AlignGetBookAnswer: answer,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishFailedGetBookAnswer(correlationID, answersRoutingkey string,
	lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: lg,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer,
		answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func isValidAlignGetRequest(request *amqp.AlignGetBookRequest) bool {
	return request != nil
}
