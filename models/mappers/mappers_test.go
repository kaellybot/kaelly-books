package mappers

import (
	"testing"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/entities"
)

func games() []amqp.Game {
	return []amqp.Game{amqp.Game_DOFUS_GAME, amqp.Game_DOFUS_TOUCH}
}

func newRequest(game amqp.Game, msgType amqp.RabbitMQMessage_Type) *amqp.RabbitMQMessage {
	return &amqp.RabbitMQMessage{
		Type:     msgType,
		Game:     game,
		Language: amqp.Language_FR,
		UserID:   "userID",
	}
}

// Every mapper must carry the request's game and language onto the answer, for
// both games, so a reply can never reach the wrong bot.
func TestMappersCarryRequestGameAndLanguage(t *testing.T) {
	t.Parallel()

	mappers := []struct {
		name        string
		requestType amqp.RabbitMQMessage_Type
		answerType  amqp.RabbitMQMessage_Type
		mapper      func(message *amqp.RabbitMQMessage) *amqp.RabbitMQMessage
		withRequest func(message *amqp.RabbitMQMessage)
	}{
		{
			name:        "job book",
			requestType: amqp.RabbitMQMessage_JOB_GET_BOOK_REQUEST,
			answerType:  amqp.RabbitMQMessage_JOB_GET_BOOK_ANSWER,
			withRequest: func(message *amqp.RabbitMQMessage) {
				message.JobGetBookRequest = &amqp.JobGetBookRequest{
					JobId: "jobID", ServerId: "serverID", Offset: 0, Size: 10,
				}
			},
			mapper: func(message *amqp.RabbitMQMessage) *amqp.RabbitMQMessage {
				return MapJobBookAnswer(message, []entities.JobBook{{UserID: "userID", Level: 100}}, 1)
			},
		},
		{
			name:        "job set",
			requestType: amqp.RabbitMQMessage_JOB_SET_REQUEST,
			answerType:  amqp.RabbitMQMessage_JOB_SET_ANSWER,
			withRequest: func(message *amqp.RabbitMQMessage) {
				message.JobSetRequest = &amqp.JobSetRequest{JobId: "jobID", ServerId: "serverID"}
			},
			mapper: MapJobSetAnswer,
		},
		{
			name:        "job user",
			requestType: amqp.RabbitMQMessage_JOB_GET_USER_REQUEST,
			answerType:  amqp.RabbitMQMessage_JOB_GET_USER_ANSWER,
			withRequest: func(message *amqp.RabbitMQMessage) {
				message.JobGetUserRequest = &amqp.JobGetUserRequest{UserId: "userID", ServerId: "serverID"}
			},
			mapper: func(message *amqp.RabbitMQMessage) *amqp.RabbitMQMessage {
				return MapJobUserAnswer(message, []entities.JobBook{{JobID: "jobID", Level: 100}})
			},
		},
		{
			name:        "align book",
			requestType: amqp.RabbitMQMessage_ALIGN_GET_BOOK_REQUEST,
			answerType:  amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER,
			withRequest: func(message *amqp.RabbitMQMessage) {
				message.AlignGetBookRequest = &amqp.AlignGetBookRequest{
					CityId: "cityID", OrderId: "orderID", ServerId: "serverID", Offset: 0, Size: 10,
				}
			},
			mapper: func(message *amqp.RabbitMQMessage) *amqp.RabbitMQMessage {
				return MapAlignBookAnswer(message, []entities.AlignmentBook{{UserID: "userID", Level: 100}}, 1)
			},
		},
		{
			name:        "align set",
			requestType: amqp.RabbitMQMessage_ALIGN_SET_REQUEST,
			answerType:  amqp.RabbitMQMessage_ALIGN_SET_ANSWER,
			withRequest: func(message *amqp.RabbitMQMessage) {
				message.AlignSetRequest = &amqp.AlignSetRequest{CityId: "cityID", OrderId: "orderID"}
			},
			mapper: MapAlignSetAnswer,
		},
		{
			name:        "align user",
			requestType: amqp.RabbitMQMessage_ALIGN_GET_USER_REQUEST,
			answerType:  amqp.RabbitMQMessage_ALIGN_GET_USER_ANSWER,
			withRequest: func(message *amqp.RabbitMQMessage) {
				message.AlignGetUserRequest = &amqp.AlignGetUserRequest{UserId: "userID", ServerId: "serverID"}
			},
			mapper: func(message *amqp.RabbitMQMessage) *amqp.RabbitMQMessage {
				return MapAlignUserAnswer(message, []entities.AlignmentBook{{CityID: "cityID", Level: 100}})
			},
		},
	}

	for _, mapper := range mappers {
		for _, game := range games() {
			t.Run(mapper.name+"/"+game.String(), func(t *testing.T) {
				t.Parallel()

				request := newRequest(game, mapper.requestType)
				mapper.withRequest(request)

				assertAnswer(t, mapper.mapper(request), game, mapper.answerType)
			})
		}
	}
}

func assertAnswer(t *testing.T, answer *amqp.RabbitMQMessage, game amqp.Game,
	answerType amqp.RabbitMQMessage_Type) {
	t.Helper()

	if answer.GetGame() != game {
		t.Errorf("Game = %v, want %v", answer.GetGame(), game)
	}
	if !answer.IsGameSet() {
		t.Error("answer game must be set, otherwise the broker refuses to publish it")
	}
	if answer.GetLanguage() != amqp.Language_FR {
		t.Errorf("Language = %v, want %v", answer.GetLanguage(), amqp.Language_FR)
	}
	if answer.GetType() != answerType {
		t.Errorf("Type = %v, want %v", answer.GetType(), answerType)
	}
	if answer.GetStatus() != amqp.RabbitMQMessage_SUCCESS {
		t.Errorf("Status = %v, want %v", answer.GetStatus(), amqp.RabbitMQMessage_SUCCESS)
	}
}
