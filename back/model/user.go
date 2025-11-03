package model

import "time"

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email" gorm:"unique"`
	Password    string    `json:"password"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Admin       bool      `json:"admin"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	HouseholdID uint      `json:"household_id" gorm:"not null"`
	Household   Household `json:"household" gorm:"foreignKey:HouseholdID;references:ID;constraint:OnDelete:CASCADE"`
}
