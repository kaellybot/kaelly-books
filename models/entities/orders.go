package entities

type Order struct {
	Id string `gorm:"primaryKey;type:varchar(100)"`
}
