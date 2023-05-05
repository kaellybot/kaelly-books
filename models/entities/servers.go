package entities

type Server struct {
	ID string `gorm:"primaryKey;type:varchar(100)"`
}
