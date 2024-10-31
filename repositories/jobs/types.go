package jobs

import (
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/utils/databases"
)

type Repository interface {
	GetBooks(jobID, serverID string, userIDs []string,
		offset, limit int) ([]entities.JobBook, int64, error)
	GetUserBook(userID, serverID string) ([]entities.JobBook, error)
	SaveUserBook(jobBook entities.JobBook) error
	DeleteUserBook(jobBook entities.JobBook) error
}

type Impl struct {
	db databases.MySQLConnection
}
