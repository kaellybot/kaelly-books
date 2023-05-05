package mappers

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/entities"
)

func MapCraftsmen(books []entities.JobBook) []*amqp.JobGetBookAnswer_Craftsman {
	craftsmen := make([]*amqp.JobGetBookAnswer_Craftsman, 0)
	for _, book := range books {
		craftsmen = append(craftsmen, &amqp.JobGetBookAnswer_Craftsman{
			UserId: book.UserID,
			Level:  book.Level,
		})
	}

	return craftsmen
}

func MapJobExperiences(books []entities.JobBook) []*amqp.JobGetUserAnswer_JobExperience {
	jobXp := make([]*amqp.JobGetUserAnswer_JobExperience, 0)
	for _, book := range books {
		jobXp = append(jobXp, &amqp.JobGetUserAnswer_JobExperience{
			JobId: book.JobID,
			Level: book.Level,
		})
	}

	return jobXp
}
