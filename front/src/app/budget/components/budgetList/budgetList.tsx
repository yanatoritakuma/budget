import { Suspense } from "react";
import { fetchBudgetList } from "@/app/api/fetchBudgetList";
import "./styles.scss";
import { formatDateForDisplay } from "@/utils/formatDateForDisplay";

export default async function BudgetList() {
  // 現在の日付を取得
  const now = new Date();
  const currentYear = now.getFullYear();
  const currentMonth = now.getMonth() + 1;

  const expenses = await fetchBudgetList({
    year: currentYear,
    month: currentMonth,
  });

  return (
    <div className="budget-list-container">
      <h2 className="budget-list-title">今月の支出一覧</h2>
      <Suspense fallback={<div className="loading-message">読み込み中...</div>}>
        {expenses.length > 0 ? (
          expenses.map((expense) => (
            <div key={expense.id} className="expense-item">
              <div className="expense-details">
                <p className="store-name">{expense.store_name}</p>
                <span className="expense-date">{formatDateForDisplay(expense.date)}</span>
              </div>
              <span className="expense-amount">
                ¥{expense.amount.toLocaleString()}
              </span>
            </div>
          ))
        ) : (
          <p className="empty-message">今月の支出はまだありません。</p>
        )}
      </Suspense>
    </div>
  );
}
