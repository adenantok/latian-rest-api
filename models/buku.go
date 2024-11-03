package models

import "gorm.io/gorm"

// Buku is a struct that represents a book
type Buku struct {
	gorm.Model
	Judul string `json:"judul" gorm:"not null" `
	Harga *int   `json:"harga" gorm:"not null" `
}

// Migrate migrates the schema
func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Buku{}); err != nil {
		panic("Migration failed: " + err.Error())
	}
}
