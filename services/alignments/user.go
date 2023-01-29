package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/models/constants"
	"github.com/kaellybot/kaelly-configurator/models/entities"
	"github.com/kaellybot/kaelly-configurator/models/mappers"
	"github.com/rs/zerolog/log"
)

func (service *AlignmentServiceImpl) UserRequest(request *amqp.AlignGetUserRequest, correlationId,
	answersRoutingkey string, lg amqp.Language) {
	if !isValidAlignGetUserRequest(request) {
		service.publishFailedGetUserAnswer(correlationId, answersRoutingkey, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationId, correlationId).
		Str(constants.LogUserId, request.UserId).
		Str(constants.LogServerId, request.ServerId).
		Msgf("Get job user request received")

	books, err := service.alignBookRepo.GetUserBook(request.UserId, request.ServerId)
	if err != nil {
		service.publishFailedGetUserAnswer(correlationId, answersRoutingkey, lg)
		return
	}

	service.publishSucceededGetUserAnswer(correlationId, request.ServerId, answersRoutingkey, books, lg)
}

func (service *AlignmentServiceImpl) publishSucceededGetUserAnswer(correlationId, serverId, answersRoutingkey string,
	books []entities.AlignmentBook, lg amqp.Language) {

	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ALIGN_GET_USER_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
		AlignGetUserAnswer: &amqp.AlignGetUserAnswer{
			Beliefs:  mappers.MapAlignExperiences(books),
			ServerId: serverId,
		},
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationId)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationId, correlationId).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *AlignmentServiceImpl) publishFailedGetUserAnswer(correlationId, answersRoutingkey string, lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ALIGN_GET_USER_ANSWER,
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

func isValidAlignGetUserRequest(request *amqp.AlignGetUserRequest) bool {
	return request != nil
}
