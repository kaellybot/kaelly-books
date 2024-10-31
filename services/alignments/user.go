package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/rs/zerolog/log"
)

func (service *Impl) UserRequest(request *amqp.AlignGetUserRequest, correlationID,
	answersRoutingkey string, lg amqp.Language) {
	if !isValidAlignGetUserRequest(request) {
		service.publishFailedGetUserAnswer(correlationID, answersRoutingkey, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Str(constants.LogUserID, request.UserId).
		Str(constants.LogServerID, request.ServerId).
		Msgf("Get job user request received")

	books, err := service.alignBookRepo.GetUserBook(request.UserId, request.ServerId)
	if err != nil {
		service.publishFailedGetUserAnswer(correlationID, answersRoutingkey, lg)
		return
	}

	answer := mappers.MapAlignUserAnswer(books, request.ServerId)
	service.publishSucceededGetUserAnswer(correlationID, answersRoutingkey, answer, lg)
}

func (service *Impl) publishSucceededGetUserAnswer(correlationID, answersRoutingkey string,
	answer *amqp.AlignGetUserAnswer, lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:               amqp.RabbitMQMessage_ALIGN_GET_USER_ANSWER,
		Status:             amqp.RabbitMQMessage_SUCCESS,
		Language:           lg,
		AlignGetUserAnswer: answer,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishFailedGetUserAnswer(correlationID, answersRoutingkey string, lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ALIGN_GET_USER_ANSWER,
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

func isValidAlignGetUserRequest(request *amqp.AlignGetUserRequest) bool {
	return request != nil
}
