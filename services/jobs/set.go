package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/models/constants"
	"github.com/kaellybot/kaelly-configurator/models/entities"
	"github.com/rs/zerolog/log"
)

func (service *JobServiceImpl) SetRequest(request *amqp.JobSetRequest, correlationId,
	answersRoutingkey string, lg amqp.Language) {

	if !isValidJobSetRequest(request) {
		service.publishFailedSetAnswer(correlationId, answersRoutingkey, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationId, correlationId).
		Str(constants.LogUserId, request.UserId).
		Str(constants.LogJobId, request.JobId).
		Str(constants.LogServerId, request.ServerId).
		Msgf("Set job request received")

	jobBook := entities.JobBook{
		UserId:   request.UserId,
		JobId:    request.JobId,
		ServerId: request.ServerId,
		Level:    request.Level,
	}

	var err error
	if request.Level > 0 {
		err = service.jobBookRepo.SaveUserBook(jobBook)
	} else {
		err = service.jobBookRepo.DeleteUserBook(jobBook)
	}
	if err != nil {
		service.publishFailedSetAnswer(correlationId, answersRoutingkey, lg)
		return
	}

	service.publishSucceededSetAnswer(correlationId, answersRoutingkey, lg)
}

func (service *JobServiceImpl) publishSucceededSetAnswer(correlationId, answersRoutingkey string, lg amqp.Language) {
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

func (service *JobServiceImpl) publishFailedSetAnswer(correlationId, answersRoutingkey string, lg amqp.Language) {
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
