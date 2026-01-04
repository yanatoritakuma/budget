package usecase

import (
	"github.com/yanatoritakuma/budget/back/domain/expense"
	"github.com/yanatoritakuma/budget/back/domain/household"
	"github.com/yanatoritakuma/budget/back/domain/user"
)

// Repositories はトランザクション内で使用されるリポジトリのセットです。
type Repositories struct {
	User      user.IUserRepository
	Household household.IHouseholdRepository
	Expense   expense.IExpenseRepository
	// 今後他のリポジトリが追加された場合は、ここに追加します。
}

// IUnitOfWork は複数のリポジトリにまたがる操作を単一のトランザクションとして
// 実行するための作業単位を定義します。
type IUnitOfWork interface {
	// Transaction は引数で受け取った関数を単一のトランザクション内で実行します。
	// 関数内でのいずれかの操作がエラーを返した場合、トランザクション全体がロールバックされます。
	// 全ての操作が成功した場合、トランザクションはコミットされます。
	Transaction(fn func(repos Repositories) error) error
}
