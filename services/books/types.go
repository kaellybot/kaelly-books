package books

import (
	"errors"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/repositories/jobs"
)

const (
	requestQueueName   = "books-requests"
	requestsRoutingkey = "requests.books"
	answersRoutingkey  = "answers.books"
)

var (
	errInvalidMessage = errors.New("Invalid books request, type is not the good one and/or the dedicated message is not filled")
)

type BooksService interface {
	Consume() error
}

type BooksServiceImpl struct {
	broker      amqp.MessageBrokerInterface
	jobBookRepo jobs.JobBookRepository
}
