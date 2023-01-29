package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/models/constants"
	alignRepo "github.com/kaellybot/kaelly-configurator/repositories/alignments"
	jobRepo "github.com/kaellybot/kaelly-configurator/repositories/jobs"
	"github.com/kaellybot/kaelly-configurator/services/alignments"
	"github.com/kaellybot/kaelly-configurator/services/books"
	"github.com/kaellybot/kaelly-configurator/services/jobs"
	"github.com/kaellybot/kaelly-configurator/utils/databases"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type ApplicationInterface interface {
	Run() error
	Shutdown()
}

type Application struct {
	booksService books.BooksService
	broker       amqp.MessageBrokerInterface
}

func New() (*Application, error) {
	// misc
	db, err := databases.New()
	if err != nil {
		return nil, err
	}

	broker, err := amqp.New(constants.RabbitMQClientId, viper.GetString(constants.RabbitMqAddress),
		[]amqp.Binding{books.GetBinding()})
	if err != nil {
		return nil, err
	}

	// repositories
	jobBooksRepo := jobRepo.New(db)
	alignBooksRepo := alignRepo.New(db)

	// services
	jobService := jobs.New(broker, jobBooksRepo)
	alignService := alignments.New(broker, alignBooksRepo)
	booksService := books.New(broker, jobService, alignService)

	return &Application{
		booksService: booksService,
		broker:       broker,
	}, nil
}

func (app *Application) Run() error {
	return app.booksService.Consume()
}

func (app *Application) Shutdown() {
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
