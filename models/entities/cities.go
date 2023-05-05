package entities

type City struct {
	ID string `gorm:"primaryKey;type:varchar(100)"`
}
