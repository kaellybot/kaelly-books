package books

import (
	"testing"

	amqp "github.com/kaellybot/kaelly-amqp"
)

// stubService records which handler was invoked, so a test can assert that an
// unsupported game never reaches the repository layer.
type stubService struct {
	called *int
}

func (stub stubService) GetBookRequest(_ amqp.Context, _ *amqp.RabbitMQMessage) { *stub.called++ }
func (stub stubService) SetRequest(_ amqp.Context, _ *amqp.RabbitMQMessage)     { *stub.called++ }
func (stub stubService) UserRequest(_ amqp.Context, _ *amqp.RabbitMQMessage)    { *stub.called++ }

type publishedReply struct {
	message       *amqp.RabbitMQMessage
	correlationID string
	replyTo       string
}

func newTestService(called *int, replies *[]publishedReply) *Impl {
	broker := amqp.Mock{
		ReplyFunc: func(msg *amqp.RabbitMQMessage, correlationID, replyTo string) error {
			*replies = append(*replies, publishedReply{
				message: msg, correlationID: correlationID, replyTo: replyTo,
			})
			return nil
		},
	}
	return New(&broker, stubService{called: called}, stubService{called: called})
}

func requestTypes() []amqp.RabbitMQMessage_Type {
	return []amqp.RabbitMQMessage_Type{
		amqp.RabbitMQMessage_JOB_GET_BOOK_REQUEST,
		amqp.RabbitMQMessage_JOB_GET_USER_REQUEST,
		amqp.RabbitMQMessage_JOB_SET_REQUEST,
		amqp.RabbitMQMessage_ALIGN_GET_BOOK_REQUEST,
		amqp.RabbitMQMessage_ALIGN_GET_USER_REQUEST,
		amqp.RabbitMQMessage_ALIGN_SET_REQUEST,
	}
}

func TestConsumeHandlesSupportedGames(t *testing.T) {
	t.Parallel()

	for _, game := range []amqp.Game{amqp.Game_DOFUS_GAME, amqp.Game_DOFUS_TOUCH} {
		for _, requestType := range requestTypes() {
			t.Run(game.String()+"/"+requestType.String(), func(t *testing.T) {
				t.Parallel()

				called := 0
				replies := make([]publishedReply, 0)
				service := newTestService(&called, &replies)

				service.consume(amqp.Context{CorrelationID: "correlationID"},
					&amqp.RabbitMQMessage{Type: requestType, Game: game})

				if called != 1 {
					t.Errorf("handler called %d time(s), want 1", called)
				}
				if len(replies) != 0 {
					t.Errorf("consume published %d reply(ies), want 0", len(replies))
				}
			})
		}
	}
}

// An unsupported game must get a FAILED reply tagged with the game it asked for,
// and must never reach the handler — otherwise it would receive another game's data.
func TestConsumeRefusesUnsupportedGame(t *testing.T) {
	t.Parallel()

	for _, requestType := range requestTypes() {
		t.Run(requestType.String(), func(t *testing.T) {
			t.Parallel()

			called := 0
			replies := make([]publishedReply, 0)
			service := newTestService(&called, &replies)

			service.consume(amqp.Context{CorrelationID: "correlationID", ReplyTo: "KaellyBot-0.answers"},
				&amqp.RabbitMQMessage{Type: requestType, Game: amqp.Game_DOFUS_RETRO})

			if called != 0 {
				t.Errorf("handler called %d time(s), want 0", called)
			}
			if len(replies) != 1 {
				t.Fatalf("consume published %d reply(ies), want 1", len(replies))
			}

			reply := replies[0]
			if reply.message.GetStatus() != amqp.RabbitMQMessage_FAILED {
				t.Errorf("Status = %v, want %v", reply.message.GetStatus(), amqp.RabbitMQMessage_FAILED)
			}
			if reply.message.GetGame() != amqp.Game_DOFUS_RETRO {
				t.Errorf("Game = %v, want %v", reply.message.GetGame(), amqp.Game_DOFUS_RETRO)
			}
			if reply.correlationID != "correlationID" {
				t.Errorf("correlationID = %q, want %q", reply.correlationID, "correlationID")
			}
			if reply.replyTo != "KaellyBot-0.answers" {
				t.Errorf("replyTo = %q, want %q", reply.replyTo, "KaellyBot-0.answers")
			}
		})
	}
}

// ANY_GAME cannot be answered: the reply would be untagged too and the broker would
// refuse it. The request is dropped and the requester falls back on its own timeout.
func TestConsumeDropsRequestWithoutGame(t *testing.T) {
	t.Parallel()

	for _, requestType := range requestTypes() {
		t.Run(requestType.String(), func(t *testing.T) {
			t.Parallel()

			called := 0
			replies := make([]publishedReply, 0)
			service := newTestService(&called, &replies)

			service.consume(amqp.Context{CorrelationID: "correlationID"},
				&amqp.RabbitMQMessage{Type: requestType})

			if called != 0 {
				t.Errorf("handler called %d time(s), want 0", called)
			}
			if len(replies) != 0 {
				t.Errorf("consume published %d reply(ies), want 0", len(replies))
			}
		})
	}
}

func TestConsumeIgnoresUnknownType(t *testing.T) {
	t.Parallel()

	called := 0
	replies := make([]publishedReply, 0)
	service := newTestService(&called, &replies)

	service.consume(amqp.Context{CorrelationID: "correlationID"},
		&amqp.RabbitMQMessage{Type: amqp.RabbitMQMessage_ABOUT_REQUEST, Game: amqp.Game_DOFUS_GAME})

	if called != 0 {
		t.Errorf("handler called %d time(s), want 0", called)
	}
	if len(replies) != 0 {
		t.Errorf("consume published %d reply(ies), want 0", len(replies))
	}
}

func TestIsGameSupported(t *testing.T) {
	t.Parallel()

	tests := map[amqp.Game]bool{
		amqp.Game_ANY_GAME:    false,
		amqp.Game_DOFUS_GAME:  true,
		amqp.Game_DOFUS_TOUCH: true,
		amqp.Game_DOFUS_RETRO: false,
	}

	for game, want := range tests {
		if got := isGameSupported(game); got != want {
			t.Errorf("isGameSupported(%v) = %v, want %v", game, got, want)
		}
	}
}
