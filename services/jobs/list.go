package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/models/mappers"
	"github.com/rs/zerolog/log"
)

func (service *Impl) GetBookRequest(request *amqp.JobGetBookRequest,
	correlationID, answersRoutingkey string, lg amqp.Language) {
	if !isValidJobGetRequest(request) {
		service.publishFailedGetBookAnswer(correlationID, answersRoutingkey, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Str(constants.LogJobID, request.JobId).
		Str(constants.LogServerID, request.ServerId).
		Msgf("Get job books request received")

	books, err := service.jobBookRepo.GetBooks(request.JobId, request.ServerId,
		request.UserIds, int(request.Limit))
	if err != nil {
		service.publishFailedGetBookAnswer(correlationID, answersRoutingkey, lg)
		return
	}

	service.publishSucceededGetBookAnswer(correlationID, request.JobId,
		request.ServerId, answersRoutingkey, books, lg)
}

func (service *Impl) publishSucceededGetBookAnswer(correlationID, jobID, serverID,
	answersRoutingkey string, books []entities.JobBook, lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_JOB_GET_BOOK_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
		JobGetBookAnswer: &amqp.JobGetBookAnswer{
			JobId:     jobID,
			ServerId:  serverID,
			Craftsmen: mappers.MapCraftsmen(books),
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
		Type:     amqp.RabbitMQMessage_JOB_GET_BOOK_ANSWER,
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

func isValidJobGetRequest(request *amqp.JobGetBookRequest) bool {
	return request != nil
}
