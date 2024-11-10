package books

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/services/alignments"
	"github.com/kaellybot/kaelly-books/services/jobs"
)

const (
	requestQueueName   = "books-requests"
	requestsRoutingkey = "requests.books"
)

type Service interface {
	Consume() error
}

type Impl struct {
	broker       amqp.MessageBroker
	jobService   jobs.Service
	alignService alignments.Service
}
