package model

import (
	"time"
)

type Budget struct {
	ID          uint      `gorm:"primaryKey"`
	HouseholdID uint      `gorm:"not null"`
	YearMonth   string    `gorm:"size:7;not null;index:idx_household_year_month,unique"` // "YYYY/MM"
	Amount      int       `gorm:"not null"`
	NotifiedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
