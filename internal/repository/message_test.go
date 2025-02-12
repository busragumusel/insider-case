package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/busragumusel/insider-case/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbUser     = "postgres"
	dbPassword = "pass"
	dbHost     = "127.0.0.1"
	dbPort     = "5432"
	dbName     = "insider_case_test"
)

func setupTestDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to PostgreSQL")
	}

	db.Exec("CREATE DATABASE insider_case_test")

	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	db.AutoMigrate(&entity.Message{})

	return db
}

func TestGetByStatus(t *testing.T) {
	db := setupTestDB()
	repo := NewMessageRepository(db)
	db.Exec("DELETE FROM messages")
	message1 := entity.Message{
		ID:          1,
		PhoneNumber: "+905551111111",
		Content:     "Oldest Message",
		Status:      "pending",
		CreatedAt:   time.Now().Add(-10 * time.Minute),
	}
	message2 := entity.Message{
		ID:          2,
		PhoneNumber: "+905552222222",
		Content:     "Newest Message",
		Status:      "pending",
		CreatedAt:   time.Now().Add(-5 * time.Minute),
	}
	db.Create(&message2)
	db.Create(&message1)
	result, err := repo.GetByStatus("pending", 10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, message1.ID, result[0].ID)
	assert.Equal(t, message2.ID, result[1].ID)
}

func TestUpdate(t *testing.T) {
	db := setupTestDB()
	repo := NewMessageRepository(db)

	db.Exec("DELETE FROM messages")

	message := entity.Message{ID: 1, PhoneNumber: "+905551111111", Content: "Test", Status: "pending"}
	db.Create(&message)

	err := repo.Update(1, "sent")

	var updatedMessage entity.Message
	db.First(&updatedMessage, 1)

	assert.NoError(t, err)
	assert.Equal(t, "sent", updatedMessage.Status)
	assert.WithinDuration(t, time.Now(), updatedMessage.SentAt, 2*time.Second)
}
