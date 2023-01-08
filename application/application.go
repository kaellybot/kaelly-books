package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-configurator/models/constants"
	"github.com/kaellybot/kaelly-configurator/repositories/jobs"
	"github.com/kaellybot/kaelly-configurator/services/books"
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
	jobBooksRepo := jobs.New(db)

	// services
	booksService, err := books.New(broker, jobBooksRepo)
	if err != nil {
		return nil, err
	}

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
