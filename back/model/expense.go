package model

import "time"

type Expense struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Amount    int       `json:"amount" gorm:"not null"`
	StoreName string    `json:"store_name" gorm:"not null"`
	Date      time.Time `json:"date" gorm:"not null"`
	Category  string    `json:"category" gorm:"not null"`
	Memo      string    `json:"memo"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

type ExpenseResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Amount    int       `json:"amount"`
	StoreName string    `json:"store_name"`
	Date      time.Time `json:"date"`
	Category  string    `json:"category"`
	Memo      string    `json:"memo"`
	CreatedAt time.Time `json:"created_at"`
}
