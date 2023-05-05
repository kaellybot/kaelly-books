package alignments

import (
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/utils/databases"
)

type Repository interface {
	GetBooks(cityID, orderID, serverID string, userIDs []string, limit int) ([]entities.AlignmentBook, error)
	GetUserBook(userID, serverID string) ([]entities.AlignmentBook, error)
	SaveUserBook(alignBook entities.AlignmentBook) error
	DeleteUserBook(alignBook entities.AlignmentBook) error
}

type Impl struct {
	db databases.MySQLConnection
}
