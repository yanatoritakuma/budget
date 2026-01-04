package repository

import (
	"github.com/yanatoritakuma/budget/back/usecase"
	"gorm.io/gorm"
)

// unitOfWork は IUnitOfWork の実装です。
type unitOfWork struct {
	db *gorm.DB
}

// NewUnitOfWork は新しい unitOfWork インスタンスを生成します。
func NewUnitOfWork(db *gorm.DB) usecase.IUnitOfWork {
	return &unitOfWork{db: db}
}

// Transaction は引数で受け取った関数を単一のトランザクション内で実行します。
func (u *unitOfWork) Transaction(fn func(usecase.Repositories) error) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		// トランザクション用の新しいリポジトリインスタンスを生成します。
		// これにより、全てのDB操作が同じトランザクション(tx)を共有します。
		repos := usecase.Repositories{
			User:      NewUserRepositoryImpl(tx),
			Household: NewHouseholdRepositoryImpl(tx),
			Expense:   NewExpenseRepositoryImpl(tx),
		}
		return fn(repos)
	})
}
