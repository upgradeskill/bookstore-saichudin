package models

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Isbn      string         `json:"isbn,omitempty"`
	Title     string         `json:"title,omitempty"`
	Author    string         `json:"author,omitempty"`
	Price     float32        `json:"price,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type BookModel struct {
	DB *gorm.DB
}

func (m BookModel) All(ctx context.Context) ([]Book, error) {
	var bks []Book
	if err := m.DB.Find(&bks).Error; err != nil {
		return nil, err
	}

	return bks, nil
}

func (m BookModel) Show(ctx context.Context, bookId uint64) (Book, error) {
	var bks Book
	if err := m.DB.First(&bks, bookId).Error; err != nil {
		return bks, err
	}

	return bks, nil
}

func (m BookModel) Create(ctx context.Context, book *Book) error {
	if err := m.DB.Create(book).Error; err != nil {
		return err
	}

	return nil
}

func (m BookModel) Update(ctx context.Context, bookId uint64, book *Book) (Book, error) {
	var bks Book
	m.DB.Model(&bks).Where("id = ?", bookId).Updates(book)

	// get udpated book for response
	if err := m.DB.First(&bks, bookId).Error; err != nil {
		return bks, err
	}

	return bks, nil
}

func (m BookModel) Delete(ctx context.Context, bookId uint64) error {
	var bks Book
	m.DB.Delete(&bks, bookId)

	return nil
}