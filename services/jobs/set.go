package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/rs/zerolog/log"
)

func (service *Impl) SetRequest(request *amqp.JobSetRequest, correlationID,
	answersRoutingkey string, lg amqp.Language) {
	if !isValidJobSetRequest(request) {
		service.publishFailedSetAnswer(correlationID, answersRoutingkey, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Str(constants.LogUserID, request.UserId).
		Str(constants.LogJobID, request.JobId).
		Str(constants.LogServerID, request.ServerId).
		Msgf("Set job request received")

	jobBook := entities.JobBook{
		UserID:   request.UserId,
		JobID:    request.JobId,
		ServerID: request.ServerId,
		Level:    request.Level,
	}

	var err error
	if request.Level > 0 {
		err = service.jobBookRepo.SaveUserBook(jobBook)
	} else {
		err = service.jobBookRepo.DeleteUserBook(jobBook)
	}
	if err != nil {
		service.publishFailedSetAnswer(correlationID, answersRoutingkey, lg)
		return
	}

	service.publishSucceededSetAnswer(correlationID, answersRoutingkey, lg)
}

func (service *Impl) publishSucceededSetAnswer(correlationID, answersRoutingkey string, lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_JOB_SET_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishFailedSetAnswer(correlationID, answersRoutingkey string, lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_JOB_SET_ANSWER,
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

func isValidJobSetRequest(request *amqp.JobSetRequest) bool {
	return request != nil
}
