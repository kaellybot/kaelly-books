package entities

type Job struct {
	Id string `gorm:"primaryKey;type:varchar(100)"`
}
