"use client";

import { useState, useEffect, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Expense } from "@/types/expense";
import { formatDateForDisplay } from "@/utils/formatDateForDisplay";
import { ButtonBox } from "@/components/elements/buttonBox/buttonBox";
import "./styles.scss";

type BudgetListComponentProps = {
  expenses: Expense[] | null;
};

function BudgetListComponent({ expenses }: BudgetListComponentProps) {
  const router = useRouter();
  const searchParams = useSearchParams();

  const getInitialDate = () => {
    const year = searchParams.get("year");
    const month = searchParams.get("month");
    if (year && month) {
      const monthIndex = parseInt(month, 10) - 1;
      return new Date(parseInt(year, 10), monthIndex, 1);
    }
    return new Date();
  };

  const [currentDate, setCurrentDate] = useState(getInitialDate());

  const year = currentDate.getFullYear();
  const month = currentDate.getMonth() + 1;

  const updateURL = (newDate: Date) => {
    const newYear = newDate.getFullYear();
    const newMonth = newDate.getMonth() + 1;
    router.push(`?year=${newYear}&month=${newMonth}`);
  };

  useEffect(() => {
    const newDate = getInitialDate();
    setCurrentDate(newDate);
  }, [searchParams]);

  const handlePrevMonth = () => {
    const newDate = new Date(
      currentDate.getFullYear(),
      currentDate.getMonth() - 1,
      1
    );
    updateURL(newDate);
  };

  const handleNextMonth = () => {
    const newDate = new Date(
      currentDate.getFullYear(),
      currentDate.getMonth() + 1,
      1
    );
    updateURL(newDate);
  };

  return (
    <div className="budget-list-container">
      <div className="budget-list-header">
        <ButtonBox onClick={handlePrevMonth} variant="outlined" size="small">
          &lt; 先月
        </ButtonBox>
        <h2 className="budget-list-title">{`${year}年${month}月 の支出一覧`}</h2>
        <ButtonBox onClick={handleNextMonth} variant="outlined" size="small">
          来月 &gt;
        </ButtonBox>
      </div>
      {expenses && expenses.length > 0 ? (
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
        <p className="empty-message">{`${year}年${month}月の支出はまだありません。`}</p>
      )}
    </div>
  );
}

type BudgetListProps = {
  expenses: Expense[] | null;
};

export default function BudgetList({ expenses }: BudgetListProps) {
  return (
    <Suspense fallback={<div className="loading-message">読み込み中...</div>}>
      <BudgetListComponent expenses={expenses} />
    </Suspense>
  );
}
