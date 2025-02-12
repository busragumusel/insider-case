package repository

import (
	"github.com/busragumusel/insider-case/internal/entity"
	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

type MessageRepo interface {
	GetByStatus(status string, limit int) ([]entity.Message, error)
	Update(id uint, status string) error
}

func NewMessageRepository(DB *gorm.DB) *MessageRepository {
	return &MessageRepository{DB}
}

func (r *MessageRepository) GetByStatus(status string, limit int) ([]entity.Message, error) {
	var messages []entity.Message

	db := r.DB

	if status != "" {
		db = db.Where("status = ?", status)
	}

	err := db.
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error

	return messages, err
}

func (r *MessageRepository) Update(id uint, status string) error {
	err := r.DB.Model(&entity.Message{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":  status,
			"sent_at": gorm.Expr("NOW()"),
		}).Error
	return err
}
