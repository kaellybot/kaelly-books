package databases

import (
	"fmt"

	"github.com/kaellybot/kaelly-books/models/constants"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConnection interface {
	GetDB() *gorm.DB
	IsConnected() bool
	Shutdown()
}

type MySQLConnectionImpl struct {
	db *gorm.DB
}

func New() (*MySQLConnectionImpl, error) {
	dbUser := viper.GetString(constants.MySQLUser)
	dbPassword := viper.GetString(constants.MySQLPassword)
	dbURL := viper.GetString(constants.MySQLURL)
	dbName := viper.GetString(constants.MySQLDatabase)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=UTC", dbUser, dbPassword, dbURL, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &MySQLConnectionImpl{db: db}, nil
}

func (connection *MySQLConnectionImpl) GetDB() *gorm.DB {
	return connection.db
}

func (connection *MySQLConnectionImpl) IsConnected() bool {
	if connection.db == nil {
		return false
	}

	dbSQL, errSQL := connection.db.DB()
	if errSQL != nil {
		return false
	}

	if errPing := dbSQL.Ping(); errPing != nil {
		return false
	}

	return true
}

func (connection *MySQLConnectionImpl) Shutdown() {
	dbSQL, err := connection.db.DB()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to shutdown database connection")
		return
	}

	if errClose := dbSQL.Close(); errClose != nil {
		log.Error().Err(errClose).Msgf("Failed to shutdown database connection")
	}
}
