package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/models/constants"
	"github.com/kaellybot/kaelly-configurator/models/entities"
	"github.com/rs/zerolog/log"
)

func (service *AlignmentServiceImpl) SetRequest(request *amqp.AlignSetRequest, correlationId,
	answersRoutingkey string, lg amqp.Language) {

	if !isValidAlignSetRequest(request) {
		service.publishFailedSetAnswer(correlationId, answersRoutingkey, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationId, correlationId).
		Str(constants.LogUserId, request.UserId).
		Str(constants.LogCityId, request.CityId).
		Str(constants.LogOrderId, request.OrderId).
		Str(constants.LogServerId, request.ServerId).
		Msgf("Set job request received")

	jobBook := entities.AlignmentBook{
		UserId:   request.UserId,
		CityId:   request.CityId,
		OrderId:  request.OrderId,
		ServerId: request.ServerId,
		Level:    request.Level,
	}

	var err error
	if request.Level > 0 {
		err = service.alignBookRepo.SaveUserBook(jobBook)
	} else {
		err = service.alignBookRepo.DeleteUserBook(jobBook)
	}
	if err != nil {
		service.publishFailedSetAnswer(correlationId, answersRoutingkey, lg)
		return
	}

	service.publishSucceededSetAnswer(correlationId, answersRoutingkey, lg)
}

func (service *AlignmentServiceImpl) publishSucceededSetAnswer(correlationId, answersRoutingkey string, lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ALIGN_SET_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationId)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationId, correlationId).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *AlignmentServiceImpl) publishFailedSetAnswer(correlationId, answersRoutingkey string, lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ALIGN_SET_ANSWER,
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

func isValidAlignSetRequest(request *amqp.AlignSetRequest) bool {
	return request != nil
}
