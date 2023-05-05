package jobs

import (
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetBooks(jobID, serverID string, userIDs []string, limit int) ([]entities.JobBook, error) {
	var jobBooks []entities.JobBook
	return jobBooks, repo.db.GetDB().
		Where("job_id = ? AND server_id = ? AND user_id IN (?)", jobID, serverID, userIDs).
		Order("level DESC").
		Limit(limit).
		Find(&jobBooks).Error
}

func (repo *Impl) GetUserBook(userID, serverID string) ([]entities.JobBook, error) {
	var jobBooks []entities.JobBook
	return jobBooks, repo.db.GetDB().
		Where("user_id = ? AND server_id = ?", userID, serverID).
		Find(&jobBooks).Error
}

func (repo *Impl) SaveUserBook(jobBook entities.JobBook) error {
	return repo.db.GetDB().Save(jobBook).Error
}

func (repo *Impl) DeleteUserBook(jobBook entities.JobBook) error {
	return repo.db.GetDB().Delete(jobBook).Error
}
