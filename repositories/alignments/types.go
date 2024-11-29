package alignments

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/utils/databases"
)

type Repository interface {
	GetBooks(cityID, orderID, serverID string, userIDs []string,
		game amqp.Game, offset, limit int) ([]entities.AlignmentBook, int64, error)
	GetUserBook(userID, serverID string, game amqp.Game) ([]entities.AlignmentBook, error)
	SaveUserBook(alignBook entities.AlignmentBook) error
	DeleteUserBook(alignBook entities.AlignmentBook) error
}

type Impl struct {
	db databases.MySQLConnection
}
