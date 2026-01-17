package user

import (
	"fmt"
	"net/mail"
	"unicode/utf8"
)

// UserID はユーザーのIDを示す値オブジェクト
type UserID uint

func (id UserID) Value() uint {
	return uint(id)
}

// Email はメールアドレスを示す値オブジェクト
type Email string

func NewEmail(address string) (*Email, error) {
	if address == "" {
		return nil, nil
	}
	_, err := mail.ParseAddress(address)
	if err != nil {
		return nil, fmt.Errorf("無効なメールアドレス形式です: %w", err)
	}
	email := Email(address)
	return &email, nil
}

func (e *Email) Value() string {
	if e == nil {
		return ""
	}
	return string(*e)
}

// Password はハッシュ化されたパスワードを示す値オブジェクト
type Password string

func NewPassword(hash string) (Password, error) {
	// ここではハッシュ化されていることを前提とするため、バリデーションは行わない
	// 必要であれば、ハッシュの形式などを検証するロジックを追加できる
	return Password(hash), nil
}

func (p Password) Value() string {
	return string(p)
}

// Name はユーザー名を示す値オブジェクト
type Name string

const MaxNameLength = 12

func NewName(name string) (Name, error) {
	if utf8.RuneCountInString(name) > MaxNameLength {
		return "", fmt.Errorf("名前は%d文字以内で入力してください", MaxNameLength)
	}
	if name == "" {
		return "", fmt.Errorf("名前は必須です")
	}
	return Name(name), nil
}

func (n Name) Value() string {
	return string(n)
}

// LineUserID はLINEのユーザーIDを示す値オブジェクト
type LineUserID string

func NewLineUserID(id string) (*LineUserID, error) {
	if id == "" {
		return nil, nil
	}
	// LINE User IDは必須ではないため、バリデーションは行わない
	lineID := LineUserID(id)
	return &lineID, nil
}

func (l *LineUserID) Value() string {
	if l == nil {
		return ""
	}
	return string(*l)
}
