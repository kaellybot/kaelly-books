package mappers

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/models/entities"
)

func MapBelievers(books []entities.AlignmentBook) []*amqp.AlignGetBookAnswer_Believer {
	believers := make([]*amqp.AlignGetBookAnswer_Believer, 0)

	for _, book := range books {
		believers = append(believers, &amqp.AlignGetBookAnswer_Believer{
			CityId:  book.CityId,
			OrderId: book.OrderId,
			UserId:  book.UserId,
			Level:   book.Level,
		})
	}

	return believers
}

func MapAlignExperiences(books []entities.AlignmentBook) []*amqp.AlignGetUserAnswer_AlignExperience {
	alignXp := make([]*amqp.AlignGetUserAnswer_AlignExperience, 0)

	for _, book := range books {
		alignXp = append(alignXp, &amqp.AlignGetUserAnswer_AlignExperience{
			CityId:  book.CityId,
			OrderId: book.OrderId,
			Level:   book.Level,
		})
	}

	return alignXp
}
