package mappers

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/entities"
)

func MapJobBookAnswer(request *amqp.JobGetBookRequest, books []entities.JobBook,
	total int64) *amqp.JobGetBookAnswer {
	craftsmen := make([]*amqp.JobGetBookAnswer_Craftsman, 0)
	for _, book := range books {
		craftsmen = append(craftsmen, &amqp.JobGetBookAnswer_Craftsman{
			UserId: book.UserID,
			Level:  book.Level,
		})
	}

	page := request.GetOffset() / request.GetSize()
	pages := int32(total) / request.GetSize()
	if int32(total)%request.GetSize() != 0 {
		pages++
	}

	return &amqp.JobGetBookAnswer{
		JobId:     request.GetJobId(),
		ServerId:  request.GetServerId(),
		Craftsmen: craftsmen,
		Page:      page,
		Pages:     pages,
		Total:     int32(total),
	}
}

func MapJobUserAnswer(books []entities.JobBook, serverID string) *amqp.JobGetUserAnswer {
	jobXp := make([]*amqp.JobGetUserAnswer_JobExperience, 0)
	for _, book := range books {
		jobXp = append(jobXp, &amqp.JobGetUserAnswer_JobExperience{
			JobId: book.JobID,
			Level: book.Level,
		})
	}

	return &amqp.JobGetUserAnswer{
		Jobs:     jobXp,
		ServerId: serverID,
	}
}
