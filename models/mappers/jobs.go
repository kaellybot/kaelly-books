package mappers

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/models/entities"
)

func MapCraftsmen(books []entities.JobBook) []*amqp.JobGetAnswer_Craftsman {
	craftsmen := make([]*amqp.JobGetAnswer_Craftsman, 0)

	for _, book := range books {
		craftsmen = append(craftsmen, &amqp.JobGetAnswer_Craftsman{
			UserId: book.UserId,
			Level:  book.Level,
		})
	}

	return craftsmen
}
