package entities

type City struct {
	Id string `gorm:"primaryKey;type:varchar(100)"`
}
