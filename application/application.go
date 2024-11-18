package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/constants"
	alignRepo "github.com/kaellybot/kaelly-books/repositories/alignments"
	jobRepo "github.com/kaellybot/kaelly-books/repositories/jobs"
	"github.com/kaellybot/kaelly-books/services/alignments"
	"github.com/kaellybot/kaelly-books/services/books"
	"github.com/kaellybot/kaelly-books/services/jobs"
	"github.com/kaellybot/kaelly-books/utils/databases"
	"github.com/kaellybot/kaelly-books/utils/insights"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() (*Impl, error) {
	// misc
	broker := amqp.New(constants.RabbitMQClientID, viper.GetString(constants.RabbitMQAddress),
		amqp.WithBindings(books.GetBinding()))
	db := databases.New()
	probes := insights.NewProbes(broker.IsConnected, db.IsConnected)
	prom := insights.NewPrometheusMetrics()

	// repositories
	jobBooksRepo := jobRepo.New(db)
	alignBooksRepo := alignRepo.New(db)

	// services
	jobService := jobs.New(broker, jobBooksRepo)
	alignService := alignments.New(broker, alignBooksRepo)
	booksService := books.New(broker, jobService, alignService)

	return &Impl{
		booksService: booksService,
		broker:       broker,
		db:           db,
		probes:       probes,
		prom:         prom,
	}, nil
}

func (app *Impl) Run() error {
	app.probes.ListenAndServe()
	app.prom.ListenAndServe()

	if err := app.db.Run(); err != nil {
		return err
	}

	if err := app.broker.Run(); err != nil {
		return err
	}

	app.booksService.Consume()
	return nil
}

func (app *Impl) Shutdown() {
	app.broker.Shutdown()
	app.db.Shutdown()
	app.prom.Shutdown()
	app.probes.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
