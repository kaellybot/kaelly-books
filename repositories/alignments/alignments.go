package alignments

import (
	"github.com/kaellybot/kaelly-configurator/models/entities"
	"github.com/kaellybot/kaelly-configurator/utils/databases"
)

type AlignmentBookRepository interface {
	GetBooks(cityId, orderId, serverId string, userIds []string, limit int) ([]entities.AlignmentBook, error)
	GetUserBook(userId, serverId string) ([]entities.AlignmentBook, error)
	SaveUserBook(alignBook entities.AlignmentBook) error
	DeleteUserBook(alignBook entities.AlignmentBook) error
}

type AlignmentBookRepositoryImpl struct {
	db databases.MySQLConnection
}

func New(db databases.MySQLConnection) *AlignmentBookRepositoryImpl {
	return &AlignmentBookRepositoryImpl{db: db}
}

func (repo *AlignmentBookRepositoryImpl) GetBooks(cityId, orderId, serverId string,
	userIds []string, limit int) ([]entities.AlignmentBook, error) {

	// TODO no optional order/city here

	var alignBooks []entities.AlignmentBook
	return alignBooks, repo.db.GetDB().
		Where("city_id = ? AND order_id = ? AND server_id = ? AND user_id IN (?)", cityId, orderId, serverId, userIds).
		Order("level DESC").
		Limit(limit).
		Find(&alignBooks).Error
}

func (repo *AlignmentBookRepositoryImpl) GetUserBook(userId, serverId string) ([]entities.AlignmentBook, error) {
	var alignBooks []entities.AlignmentBook
	return alignBooks, repo.db.GetDB().
		Where("user_id = ? AND server_id = ?", userId, serverId).
		Find(&alignBooks).Error
}

func (repo *AlignmentBookRepositoryImpl) SaveUserBook(alignBook entities.AlignmentBook) error {
	return repo.db.GetDB().Save(alignBook).Error
}

func (repo *AlignmentBookRepositoryImpl) DeleteUserBook(alignBook entities.AlignmentBook) error {
	return repo.db.GetDB().Delete(alignBook).Error
}
