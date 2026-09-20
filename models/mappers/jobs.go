package mappers

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/entities"
)

func MapJobBookAnswer(message *amqp.RabbitMQMessage, books []entities.JobBook,
	total int64) *amqp.RabbitMQMessage {
	request := message.GetJobGetBookRequest()
	craftsmen := make([]*amqp.JobGetBookAnswer_Craftsman, 0)
	for _, book := range books {
		craftsmen = append(craftsmen, &amqp.JobGetBookAnswer_Craftsman{
			UserId: book.UserID,
			Level:  book.Level,
		})
	}

	page := request.GetOffset() / request.GetSize()
	pages := total / request.GetSize()
	if total%request.GetSize() != 0 {
		pages++
	}

	answer := amqp.NewReply(message, amqp.RabbitMQMessage_JOB_GET_BOOK_ANSWER,
		amqp.RabbitMQMessage_SUCCESS)
	answer.JobGetBookAnswer = &amqp.JobGetBookAnswer{
		JobId:     request.GetJobId(),
		ServerId:  request.GetServerId(),
		Craftsmen: craftsmen,
		Page:      page,
		Pages:     pages,
		Total:     total,
	}
	return answer
}

func MapJobSetAnswer(message *amqp.RabbitMQMessage) *amqp.RabbitMQMessage {
	return amqp.NewReply(message, amqp.RabbitMQMessage_JOB_SET_ANSWER,
		amqp.RabbitMQMessage_SUCCESS)
}

func MapJobUserAnswer(message *amqp.RabbitMQMessage,
	books []entities.JobBook) *amqp.RabbitMQMessage {
	jobXp := make([]*amqp.JobGetUserAnswer_JobExperience, 0)
	for _, book := range books {
		jobXp = append(jobXp, &amqp.JobGetUserAnswer_JobExperience{
			JobId: book.JobID,
			Level: book.Level,
		})
	}

	answer := amqp.NewReply(message, amqp.RabbitMQMessage_JOB_GET_USER_ANSWER,
		amqp.RabbitMQMessage_SUCCESS)
	answer.JobGetUserAnswer = &amqp.JobGetUserAnswer{
		Jobs:     jobXp,
		ServerId: message.GetJobGetUserRequest().GetServerId(),
	}
	return answer
}
