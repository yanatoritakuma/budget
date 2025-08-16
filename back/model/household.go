package model

import "time"

type Household struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"not null"`
	InviteCode string    `json:"invite_code" gorm:"unique"`
	Users      []User    `json:"users"` // A household has many users
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}