package household

import (
	"fmt"
	"unicode/utf8"
)

// HouseholdID は家計のIDを示す値オブジェクト
type HouseholdID uint

func (id HouseholdID) Value() uint {
	return uint(id)
}

// Name は家計名を示す値オブジェクト
type Name string

func NewName(name string) (Name, error) {
	if name == "" {
		return "", fmt.Errorf("家計名は必須です")
	}
	// Note: You can add more validation rules here if needed, e.g., length.
	return Name(name), nil
}

func (n Name) Value() string {
	return string(n)
}

// InviteCode は招待コードを示す値オブジェクト
type InviteCode string

const InviteCodeLength = 16

func NewInviteCode(code string) (InviteCode, error) {
	if utf8.RuneCountInString(code) != InviteCodeLength {
		return "", fmt.Errorf("招待コードは%d文字である必要があります", InviteCodeLength)
	}
	return InviteCode(code), nil
}

func (c InviteCode) Value() string {
	return string(c)
}
