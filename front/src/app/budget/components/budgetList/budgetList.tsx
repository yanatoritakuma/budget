"use client";

import { useState, useEffect, useCallback, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { fetchBudgetList } from "@/app/api/fetchBudgetList";
import { Expense } from "@/types/expense";
import { formatDateForDisplay } from "@/utils/formatDateForDisplay";
import { ButtonBox } from "@/components/elements/buttonBox/buttonBox";
import "./styles.scss";

function BudgetListComponent() {
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
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const year = currentDate.getFullYear();
  const month = currentDate.getMonth() + 1;

  const updateURL = (newDate: Date) => {
    const newYear = newDate.getFullYear();
    const newMonth = newDate.getMonth() + 1;
    // Use router.push to update the URL without a full page reload
    router.push(`?year=${newYear}&month=${newMonth}`);
  };

  const fetchExpenses = useCallback(async () => {
    setIsLoading(true);
    try {
      const fetchedExpenses = await fetchBudgetList({ year, month });
      setExpenses(fetchedExpenses || []);
    } catch (error) {
      console.error("Failed to fetch expenses:", error);
      setExpenses([]);
    } finally {
      setIsLoading(false);
    }
  }, [year, month]);

  useEffect(() => {
    // This effect syncs the component state with the URL's query params.
    // It runs when the component mounts and whenever the searchParams change.
    const newDate = getInitialDate();
    setCurrentDate(newDate);
  }, [searchParams]);

  useEffect(() => {
    // This effect fetches the data whenever the date changes.
    fetchExpenses();
  }, [currentDate, fetchExpenses]);

  const handlePrevMonth = () => {
    const newDate = new Date(currentDate.getFullYear(), currentDate.getMonth() - 1, 1);
    updateURL(newDate);
  };

  const handleNextMonth = () => {
    const newDate = new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 1);
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
      {isLoading ? (
        <div className="loading-message">読み込み中...</div>
      ) : expenses.length > 0 ? (
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

// Since searchParams can only be used in a Client Component,
// and we want to use it with Suspense, we wrap it like this.
export default function BudgetList() {
  return (
    <Suspense fallback={<div className="loading-message">読み込み中...</div>}>
      <BudgetListComponent />
    </Suspense>
  );
}
