package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/entities"
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

	books, err := service.alignBookRepo.GetBooks(request.CityId, request.OrderId,
		request.ServerId, request.UserIds, int(request.Limit))
	if err != nil {
		service.publishFailedGetBookAnswer(correlationID, answersRoutingkey, lg)
		return
	}

	service.publishSucceededGetBookAnswer(correlationID, request.ServerId, answersRoutingkey, books, lg)
}

func (service *Impl) publishSucceededGetBookAnswer(correlationID, serverID,
	answersRoutingkey string, books []entities.AlignmentBook, lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
		AlignGetBookAnswer: &amqp.AlignGetBookAnswer{
			ServerId:  serverID,
			Believers: mappers.MapBelievers(books),
		},
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
