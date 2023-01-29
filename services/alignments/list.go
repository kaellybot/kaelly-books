package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/models/constants"
	"github.com/kaellybot/kaelly-configurator/models/entities"
	"github.com/kaellybot/kaelly-configurator/models/mappers"
	"github.com/rs/zerolog/log"
)

func (service *AlignmentServiceImpl) GetBookRequest(request *amqp.AlignGetBookRequest,
	correlationId, answersRoutingkey string, lg amqp.Language) {

	if !isValidAlignGetRequest(request) {
		service.publishFailedGetBookAnswer(correlationId, answersRoutingkey, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationId, correlationId).
		Str(constants.LogCityId, request.CityId).
		Str(constants.LogOrderId, request.OrderId).
		Str(constants.LogServerId, request.ServerId).
		Msgf("Get job books request received")

	books, err := service.alignBookRepo.GetBooks(request.CityId, request.OrderId,
		request.ServerId, request.UserIds, int(request.Limit))
	if err != nil {
		service.publishFailedGetBookAnswer(correlationId, answersRoutingkey, lg)
		return
	}

	service.publishSucceededGetBookAnswer(correlationId, request.ServerId, answersRoutingkey, books, lg)
}

func (service *AlignmentServiceImpl) publishSucceededGetBookAnswer(correlationId, serverId,
	answersRoutingkey string, books []entities.AlignmentBook, lg amqp.Language) {

	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
		AlignGetBookAnswer: &amqp.AlignGetBookAnswer{
			ServerId:  serverId,
			Believers: mappers.MapBelievers(books),
		},
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationId)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationId, correlationId).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *AlignmentServiceImpl) publishFailedGetBookAnswer(correlationId, answersRoutingkey string,
	lg amqp.Language) {

	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: lg,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer,
		answersRoutingkey, correlationId)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationId, correlationId).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func isValidAlignGetRequest(request *amqp.AlignGetBookRequest) bool {
	return request != nil
}
