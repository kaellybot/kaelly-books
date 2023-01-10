package jobs

import (
	"github.com/kaellybot/kaelly-configurator/models/entities"
	"github.com/kaellybot/kaelly-configurator/utils/databases"
	"gorm.io/gorm/clause"
)

type JobBookRepository interface {
	GetBooks(jobId, serverId string, userIds []string, limit int) ([]entities.JobBook, error)
	GetUserBook(userId, serverId string) ([]entities.JobBook, error)
	SaveUserBook(jobBook entities.JobBook) error
}

type JobBookRepositoryImpl struct {
	db databases.MySQLConnection
}

func New(db databases.MySQLConnection) *JobBookRepositoryImpl {
	return &JobBookRepositoryImpl{db: db}
}

func (repo *JobBookRepositoryImpl) GetBooks(jobId, serverId string,
	userIds []string, limit int) ([]entities.JobBook, error) {

	var jobBooks []entities.JobBook
	return jobBooks, repo.db.GetDB().
		Where("job_id = ? AND server_id = ? AND user_id IN (?)", jobId, serverId, userIds).
		Order("level DESC").
		Limit(limit).
		Find(&jobBooks).Error
}

func (repo *JobBookRepositoryImpl) GetUserBook(userId, serverId string) ([]entities.JobBook, error) {
	var jobBooks []entities.JobBook
	return jobBooks, repo.db.GetDB().
		Where("user_id = ? AND server_id = ?", userId, serverId).
		Find(&jobBooks).Error
}

func (repo *JobBookRepositoryImpl) SaveUserBook(jobBook entities.JobBook) error {
	return repo.db.GetDB().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&jobBook).Error
}
