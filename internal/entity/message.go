package entity

import "time"

const (
	StatusPending = "pending"
	StatusSent    = "sent"
)

type Message struct {
	ID          uint      `gorm:"primaryKey"`
	PhoneNumber string    `gorm:"size:20;not null"`
	Content     string    `gorm:"size:160;not null"`
	Status      string    `gorm:"size:10;default:pending"`
	CreatedAt   time.Time `gorm:"default:null"`
	SentAt      time.Time `gorm:"default:null"`
}
