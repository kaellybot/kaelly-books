package mappers

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/entities"
)

func MapAlignBookAnswer(request *amqp.AlignGetBookRequest, books []entities.AlignmentBook,
	total int64) *amqp.AlignGetBookAnswer {
	believers := make([]*amqp.AlignGetBookAnswer_Believer, 0)
	for _, book := range books {
		believers = append(believers, &amqp.AlignGetBookAnswer_Believer{
			CityId:  book.CityID,
			OrderId: book.OrderID,
			UserId:  book.UserID,
			Level:   book.Level,
		})
	}

	page := request.GetOffset() / request.GetSize()
	pages := int32(total) / request.GetSize()
	if int32(total)%request.GetSize() != 0 {
		pages++
	}

	return &amqp.AlignGetBookAnswer{
		CityId:    request.GetCityId(),
		OrderId:   request.GetOrderId(),
		ServerId:  request.GetServerId(),
		Believers: believers,
		Page:      page,
		Pages:     pages,
		Total:     int32(total),
	}
}

func MapAlignUserAnswer(books []entities.AlignmentBook,
	serverID string) *amqp.AlignGetUserAnswer {
	alignXp := make([]*amqp.AlignGetUserAnswer_AlignExperience, 0)
	for _, book := range books {
		alignXp = append(alignXp, &amqp.AlignGetUserAnswer_AlignExperience{
			CityId:  book.CityID,
			OrderId: book.OrderID,
			Level:   book.Level,
		})
	}

	return &amqp.AlignGetUserAnswer{
		Beliefs:  alignXp,
		ServerId: serverID,
	}
}
