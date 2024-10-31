package alignments

import (
	"github.com/kaellybot/kaelly-books/models/entities"
	"github.com/kaellybot/kaelly-books/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetBooks(cityID, orderID, serverID string, userIDs []string,
	offset, limit int) ([]entities.AlignmentBook, int64, error) {
	var total int64
	var alignBooks []entities.AlignmentBook
	baseQuery := repo.db.GetDB().
		Where("server_id = ? AND user_id IN (?)", serverID, userIDs)

	if len(cityID) > 0 {
		baseQuery = baseQuery.Where("city_id = ?", cityID)
	}

	if len(orderID) > 0 {
		baseQuery = baseQuery.Where("order_id = ?", orderID)
	}

	if errTotal := baseQuery.Model(&entities.AlignmentBook{}).
		Count(&total).Error; errTotal != nil {
		return nil, 0, errTotal
	}

	return alignBooks, total, baseQuery.
		Order("level DESC").
		Offset(offset).
		Limit(limit).
		Find(&alignBooks).Error
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
