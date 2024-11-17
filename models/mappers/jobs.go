package mappers

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/entities"
)

func MapJobBookAnswer(request *amqp.JobGetBookRequest, books []entities.JobBook,
	total int64, lg amqp.Language) *amqp.RabbitMQMessage {
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

	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_JOB_GET_BOOK_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
		JobGetBookAnswer: &amqp.JobGetBookAnswer{
			JobId:     request.GetJobId(),
			ServerId:  request.GetServerId(),
			Craftsmen: craftsmen,
			Page:      page,
			Pages:     pages,
			Total:     total,
		},
	}
}

func MapJobSetAnswer(lg amqp.Language) *amqp.RabbitMQMessage {
	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_JOB_SET_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
	}
}

func MapJobUserAnswer(books []entities.JobBook, serverID string,
	lg amqp.Language) *amqp.RabbitMQMessage {
	jobXp := make([]*amqp.JobGetUserAnswer_JobExperience, 0)
	for _, book := range books {
		jobXp = append(jobXp, &amqp.JobGetUserAnswer_JobExperience{
			JobId: book.JobID,
			Level: book.Level,
		})
	}

	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_JOB_GET_USER_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
		JobGetUserAnswer: &amqp.JobGetUserAnswer{
			Jobs:     jobXp,
			ServerId: serverID,
		},
	}
}
