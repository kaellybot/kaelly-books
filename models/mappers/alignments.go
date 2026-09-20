package mappers

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/entities"
)

func MapAlignBookAnswer(message *amqp.RabbitMQMessage, books []entities.AlignmentBook,
	total int64) *amqp.RabbitMQMessage {
	request := message.GetAlignGetBookRequest()
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
	pages := total / request.GetSize()
	if total%request.GetSize() != 0 {
		pages++
	}

	answer := amqp.NewReply(message, amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER,
		amqp.RabbitMQMessage_SUCCESS)
	answer.AlignGetBookAnswer = &amqp.AlignGetBookAnswer{
		CityId:    request.GetCityId(),
		OrderId:   request.GetOrderId(),
		ServerId:  request.GetServerId(),
		Believers: believers,
		Page:      page,
		Pages:     pages,
		Total:     total,
	}
	return answer
}

func MapAlignSetAnswer(message *amqp.RabbitMQMessage) *amqp.RabbitMQMessage {
	return amqp.NewReply(message, amqp.RabbitMQMessage_ALIGN_SET_ANSWER,
		amqp.RabbitMQMessage_SUCCESS)
}

func MapAlignUserAnswer(message *amqp.RabbitMQMessage,
	books []entities.AlignmentBook) *amqp.RabbitMQMessage {
	alignXp := make([]*amqp.AlignGetUserAnswer_AlignExperience, 0)
	for _, book := range books {
		alignXp = append(alignXp, &amqp.AlignGetUserAnswer_AlignExperience{
			CityId:  book.CityID,
			OrderId: book.OrderID,
			Level:   book.Level,
		})
	}

	answer := amqp.NewReply(message, amqp.RabbitMQMessage_ALIGN_GET_USER_ANSWER,
		amqp.RabbitMQMessage_SUCCESS)
	answer.AlignGetUserAnswer = &amqp.AlignGetUserAnswer{
		Beliefs:  alignXp,
		ServerId: message.GetAlignGetUserRequest().GetServerId(),
	}
	return answer
}
