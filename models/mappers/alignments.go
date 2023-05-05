package mappers

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/entities"
)

func MapBelievers(books []entities.AlignmentBook) []*amqp.AlignGetBookAnswer_Believer {
	believers := make([]*amqp.AlignGetBookAnswer_Believer, 0)
	for _, book := range books {
		believers = append(believers, &amqp.AlignGetBookAnswer_Believer{
			CityId:  book.CityID,
			OrderId: book.OrderID,
			UserId:  book.UserID,
			Level:   book.Level,
		})
	}

	return believers
}

func MapAlignExperiences(books []entities.AlignmentBook) []*amqp.AlignGetUserAnswer_AlignExperience {
	alignXp := make([]*amqp.AlignGetUserAnswer_AlignExperience, 0)
	for _, book := range books {
		alignXp = append(alignXp, &amqp.AlignGetUserAnswer_AlignExperience{
			CityId:  book.CityID,
			OrderId: book.OrderID,
			Level:   book.Level,
		})
	}

	return alignXp
}
