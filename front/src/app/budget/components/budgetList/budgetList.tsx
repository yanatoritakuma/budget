import { Suspense } from "react";
import { fetchBudgetList } from "@/app/api/fetchBudgetList";
import "./styles.scss";
import { formatDateForDisplay } from "@/utils/formatDateForDisplay";
import { ButtonBox } from "@/components/elements/buttonBox/buttonBox";

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
      <div className="budget-list-header">
        <ButtonBox variant="outlined" size="small">
          &lt; 先月
        </ButtonBox>
        <h2 className="budget-list-title">{`${currentYear}年${currentMonth}月 の支出一覧`}</h2>
        <ButtonBox variant="outlined" size="small">
          来月 &gt;
        </ButtonBox>
      </div>
      <Suspense fallback={<div className="loading-message">読み込み中...</div>}>
        {expenses?.length > 0 ? (
          expenses.map((expense) => (
            <div key={expense.id} className="expense-item">
              <div className="expense-details">
                <p className="store-name">{expense.store_name}</p>
                <span className="expense-date">
                  {formatDateForDisplay(expense.date)}
                </span>
                {expense.memo && (
                  <p className="expense-memo">メモ: {expense.memo}</p>
                )}
                {expense.payer_name && (
                  <p className="expense-payer">支払者: {expense.payer_name}</p>
                )}
              </div>
              <span className="expense-amount">
                ¥{expense.amount.toLocaleString()}
              </span>
            </div>
          ))
        ) : (
          <p className="empty-message">{`${currentYear}年${currentMonth}月の支出はまだありません。`}</p>
        )}
      </Suspense>
    </div>
  );
}
