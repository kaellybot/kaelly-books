package entities

type Server struct {
	Id string `gorm:"primaryKey;type:varchar(100)"`
}
