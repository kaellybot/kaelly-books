package mappers

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/models/entities"
)

func MapCraftsmen(books []entities.JobBook) []*amqp.JobGetBookAnswer_Craftsman {
	craftsmen := make([]*amqp.JobGetBookAnswer_Craftsman, 0)

	for _, book := range books {
		craftsmen = append(craftsmen, &amqp.JobGetBookAnswer_Craftsman{
			UserId: book.UserId,
			Level:  book.Level,
		})
	}

	return craftsmen
}

func MapJobExperiences(books []entities.JobBook) []*amqp.JobGetUserAnswer_JobExperience {
	jobXp := make([]*amqp.JobGetUserAnswer_JobExperience, 0)

	for _, book := range books {
		jobXp = append(jobXp, &amqp.JobGetUserAnswer_JobExperience{
			JobId: book.JobId,
			Level: book.Level,
		})
	}

	return jobXp
}
