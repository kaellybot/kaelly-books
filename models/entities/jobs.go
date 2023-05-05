package entities

type Job struct {
	ID string `gorm:"primaryKey;type:varchar(100)"`
}
