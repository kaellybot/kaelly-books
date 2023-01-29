package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/models/constants"
	"github.com/kaellybot/kaelly-configurator/models/entities"
	"github.com/kaellybot/kaelly-configurator/models/mappers"
	"github.com/rs/zerolog/log"
)

func (service *JobServiceImpl) UserRequest(request *amqp.JobGetUserRequest, correlationId,
	answersRoutingkey string, lg amqp.Language) {
	if !isValidJobGetUserRequest(request) {
		service.publishFailedGetUserAnswer(correlationId, answersRoutingkey, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationId, correlationId).
		Str(constants.LogUserId, request.UserId).
		Str(constants.LogServerId, request.ServerId).
		Msgf("Get job user request received")

	books, err := service.jobBookRepo.GetUserBook(request.UserId, request.ServerId)
	if err != nil {
		service.publishFailedGetUserAnswer(correlationId, answersRoutingkey, lg)
		return
	}

	service.publishSucceededGetUserAnswer(correlationId, request.ServerId, answersRoutingkey, books, lg)
}

func (service *JobServiceImpl) publishSucceededGetUserAnswer(correlationId, serverId, answersRoutingkey string,
	books []entities.JobBook, lg amqp.Language) {

	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_JOB_GET_USER_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
		JobGetUserAnswer: &amqp.JobGetUserAnswer{
			Jobs:     mappers.MapJobExperiences(books),
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

func (service *JobServiceImpl) publishFailedGetUserAnswer(correlationId, answersRoutingkey string, lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_JOB_GET_USER_ANSWER,
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

func isValidJobGetUserRequest(request *amqp.JobGetUserRequest) bool {
	return request != nil
}
