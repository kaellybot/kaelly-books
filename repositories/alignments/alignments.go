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

	var alignBooks []entities.AlignmentBook
	query := repo.db.GetDB().
		Where("server_id = ? AND user_id IN (?)", serverId, userIds).
		Order("level DESC").
		Limit(limit)

	if len(cityId) > 0 {
		query = query.Where("city_id = ?", cityId)
	}

	if len(orderId) > 0 {
		query = query.Where("order_id = ?", orderId)
	}

	return alignBooks, query.Find(&alignBooks).Error
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
