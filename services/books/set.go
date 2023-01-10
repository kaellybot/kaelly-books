package books

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/models/constants"
	"github.com/kaellybot/kaelly-configurator/models/entities"
	"github.com/rs/zerolog/log"
)

func (service *BooksServiceImpl) setRequest(message *amqp.RabbitMQMessage, correlationId string) {
	request := message.JobSetRequest
	if !isValidJobSetRequest(request) {
		service.publishFailedSetAnswer(correlationId, message.Language)
		return
	}

	log.Info().Str(constants.LogCorrelationId, correlationId).
		Msgf("Set job request received")

	err := service.jobBookRepo.SaveUserBook(entities.JobBook{
		UserId:   request.UserId,
		JobId:    request.JobId,
		ServerId: request.ServerId,
		Level:    request.Level,
	})
	if err != nil {
		service.publishFailedGetAnswer(correlationId, message.Language)
		return
	}

	service.publishSucceededSetAnswer(correlationId, message.Language)
}

func (service *BooksServiceImpl) publishSucceededSetAnswer(correlationId string,
	lg amqp.RabbitMQMessage_Language) {

	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_JOB_SET_ANSWER,
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

func (service *BooksServiceImpl) publishFailedSetAnswer(correlationId string,
	lg amqp.RabbitMQMessage_Language) {

	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_JOB_SET_ANSWER,
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

func isValidJobSetRequest(request *amqp.JobSetRequest) bool {
	return request != nil
}
