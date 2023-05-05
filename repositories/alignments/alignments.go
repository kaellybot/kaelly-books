package alignments

import (
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetBooks(cityID, orderID, serverID string,
	userIDs []string, limit int) ([]entities.AlignmentBook, error) {
	var alignBooks []entities.AlignmentBook
	query := repo.db.GetDB().
		Where("server_id = ? AND user_id IN (?)", serverID, userIDs).
		Order("level DESC").
		Limit(limit)

	if len(cityID) > 0 {
		query = query.Where("city_id = ?", cityID)
	}

	if len(orderID) > 0 {
		query = query.Where("order_id = ?", orderID)
	}

	return alignBooks, query.Find(&alignBooks).Error
}

func (repo *Impl) GetUserBook(userID, serverID string) ([]entities.AlignmentBook, error) {
	var alignBooks []entities.AlignmentBook
	return alignBooks, repo.db.GetDB().
		Where("user_id = ? AND server_id = ?", userID, serverID).
		Find(&alignBooks).Error
}

func (repo *Impl) SaveUserBook(alignBook entities.AlignmentBook) error {
	return repo.db.GetDB().Save(alignBook).Error
}

func (repo *Impl) DeleteUserBook(alignBook entities.AlignmentBook) error {
	return repo.db.GetDB().Delete(alignBook).Error
}
