package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/rs/zerolog/log"
)

func (service *Impl) UserRequest(request *amqp.JobGetUserRequest, correlationID,
	answersRoutingkey string, lg amqp.Language) {
	if !isValidJobGetUserRequest(request) {
		service.publishFailedGetUserAnswer(correlationID, answersRoutingkey, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Str(constants.LogUserID, request.UserId).
		Str(constants.LogServerID, request.ServerId).
		Msgf("Get job user request received")

	books, err := service.jobBookRepo.GetUserBook(request.UserId, request.ServerId)
	if err != nil {
		service.publishFailedGetUserAnswer(correlationID, answersRoutingkey, lg)
		return
	}

	service.publishSucceededGetUserAnswer(correlationID, request.ServerId, answersRoutingkey, books, lg)
}

func (service *Impl) publishSucceededGetUserAnswer(correlationID, serverID, answersRoutingkey string,
	books []entities.JobBook, lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_JOB_GET_USER_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
		JobGetUserAnswer: &amqp.JobGetUserAnswer{
			Jobs:     mappers.MapJobExperiences(books),
			ServerId: serverID,
		},
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
		Type:     amqp.RabbitMQMessage_JOB_GET_USER_ANSWER,
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

func isValidJobGetUserRequest(request *amqp.JobGetUserRequest) bool {
	return request != nil
}
