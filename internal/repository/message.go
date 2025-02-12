package repository

import (
	"context"
	"github.com/busragumusel/insider-case/internal/entity"
	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

type MessageRepo interface {
	GetByStatus(ctx context.Context, status string, limit int) ([]entity.Message, error)
	Update(ctx context.Context, id uint, status string) error
}

func NewMessageRepository(DB *gorm.DB) *MessageRepository {
	return &MessageRepository{DB}
}

func (r *MessageRepository) GetByStatus(
	ctx context.Context,
	status string,
	limit int,
) ([]entity.Message, error) {
	var messages []entity.Message

	db := r.DB.WithContext(ctx)

	if status != "" {
		db = db.Where("status = ?", status)
	}

	err := db.
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error

	return messages, err
}

func (r *MessageRepository) Update(ctx context.Context, id uint, status string) error {
	err := r.DB.WithContext(ctx).
		Model(&entity.Message{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":  status,
			"sent_at": gorm.Expr("NOW()"),
		}).Error
	return err
}
