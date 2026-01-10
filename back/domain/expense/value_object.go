package expense

import (
	"fmt"
	"unicode/utf8"

	"github.com/yanatoritakuma/budget/back/domain/user"
)

// ExpenseID は支出のIDを示す値オブジェクト
type ExpenseID uint

func (id ExpenseID) Value() uint {
	return uint(id)
}

// Amount は金額を示す値オブジェクト
type Amount int

const MinAmount = 0

func NewAmount(amount int) (Amount, error) {
	if amount <= MinAmount {
		return 0, fmt.Errorf("金額は%dより大きい値を入力してください", MinAmount)
	}
	return Amount(amount), nil
}

func (a Amount) Value() int {
	return int(a)
}

// StoreName は店名を示す値オブジェクト
type StoreName string

const MaxStoreNameLength = 255

func NewStoreName(name string) (StoreName, error) {
	if utf8.RuneCountInString(name) > MaxStoreNameLength {
		return "", fmt.Errorf("店名は%d文字以内で入力してください", MaxStoreNameLength)
	}
	return StoreName(name), nil
}

func (n StoreName) Value() string {
	return string(n)
}

// Category はカテゴリを示す値オブジェクト
type Category string

func NewCategory(category string) (Category, error) {
	// TODO: カテゴリのバリデーションルール（例: 許容リスト）を追加
	if category == "" {
		return "", fmt.Errorf("カテゴリは必須です")
	}
	return Category(category), nil
}

func (c Category) Value() string {
	return string(c)
}

// Memo はメモを示す値オブジェクト
type Memo string

const MaxMemoLength = 1000

func NewMemo(memo string) (Memo, error) {
	if utf8.RuneCountInString(memo) > MaxMemoLength {
		return "", fmt.Errorf("メモは%d文字以内で入力してください", MaxMemoLength)
	}
	return Memo(memo), nil
}

func (m Memo) Value() string {
	return string(m)
}

// UserID はユーザーIDの値オブジェクト
type UserID user.UserID

// PayerID は支払者IDの値オブジェクト
type PayerID user.UserID
