package entities

import amqp "github.com/kaellybot/kaelly-amqp"

type City struct {
	ID   string    `gorm:"primaryKey;type:varchar(100)"`
	Game amqp.Game `gorm:"primaryKey"`
}
