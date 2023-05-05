package entities

type Order struct {
	ID string `gorm:"primaryKey;type:varchar(100)"`
}
