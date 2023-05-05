package jobs

import (
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/utils/databases"
)

type Repository interface {
	GetBooks(jobID, serverID string, userIDs []string, limit int) ([]entities.JobBook, error)
	GetUserBook(userID, serverID string) ([]entities.JobBook, error)
	SaveUserBook(jobBook entities.JobBook) error
	DeleteUserBook(jobBook entities.JobBook) error
}

type Impl struct {
	db databases.MySQLConnection
}
