package entities

type JobBook struct {
	UserID   string `gorm:"primaryKey;type:varchar(100);"`
	JobID    string `gorm:"primaryKey;type:varchar(100);"`
	ServerID string `gorm:"primaryKey;type:varchar(100);"`
	Level    int64
	Job      Job    `gorm:"foreignKey:JobID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Server   Server `gorm:"foreignKey:ServerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type AlignmentBook struct {
	UserID   string `gorm:"primaryKey;type:varchar(100);"`
	CityID   string `gorm:"primaryKey;type:varchar(100);"`
	OrderID  string `gorm:"primaryKey;type:varchar(100);"`
	ServerID string `gorm:"primaryKey;type:varchar(100);"`
	Level    int64
	City     City   `gorm:"foreignKey:CityID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Order    Order  `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Server   Server `gorm:"foreignKey:ServerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
