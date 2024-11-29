package jobs

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetBooks(jobID, serverID string, userIDs []string,
	game amqp.Game, offset, limit int) ([]entities.JobBook, int64, error) {
	var total int64
	var jobBooks []entities.JobBook
	baseQuery := repo.db.GetDB().
		Where("job_id = ? AND server_id = ? AND game = ? AND user_id IN (?)", jobID, serverID, game, userIDs)

	if errTotal := baseQuery.Model(&entities.JobBook{}).
		Count(&total).Error; errTotal != nil {
		return nil, 0, errTotal
	}

	return jobBooks, total, baseQuery.
		Order("level DESC").
		Offset(offset).
		Limit(limit).
		Find(&jobBooks).Error
}

func (repo *Impl) GetUserBook(userID, serverID string, game amqp.Game) ([]entities.JobBook, error) {
	var jobBooks []entities.JobBook
	return jobBooks, repo.db.GetDB().
		Where("user_id = ? AND server_id = ? AND game = ?", userID, serverID, game).
		Find(&jobBooks).Error
}

func (repo *Impl) SaveUserBook(jobBook entities.JobBook) error {
	return repo.db.GetDB().Save(jobBook).Error
}

func (repo *Impl) DeleteUserBook(jobBook entities.JobBook) error {
	return repo.db.GetDB().Delete(jobBook).Error
}
