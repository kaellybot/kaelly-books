package entities

type JobBook struct {
	UserId   string `gorm:"primaryKey;type:varchar(100);"`
	JobId    string `gorm:"primaryKey;type:varchar(100);"`
	ServerId string `gorm:"primaryKey;type:varchar(100);"`
	Level    int64
	Job      Job    `gorm:"foreignKey:JobId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Server   Server `gorm:"foreignKey:ServerId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
