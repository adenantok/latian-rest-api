package models

import (
	"time"

	"gorm.io/gorm"
)

// Buku is a struct that represents a book
type Buku struct {
	Id        *int      `json:"id" gorm:"primaryKey;not null" `
	Judul     string    `json:"judul" gorm:"not null" `
	Harga     *int      `json:"harga" gorm:"not null" `
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Migrate migrates the schema
func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Buku{}); err != nil {
		panic("Migration failed: " + err.Error())
	}
}
