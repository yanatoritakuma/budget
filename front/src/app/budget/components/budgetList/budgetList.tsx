import { Suspense } from "react";
import { fetchBudgetList } from "@/app/api/fetchBudgetList";
import "@/app/budget/components/budgetList/styles.scss";

export default async function BudgetList() {
  // 現在の日付を取得
  const now = new Date();
  const currentYear = now.getFullYear();
  const currentMonth = now.getMonth() + 1;

  const expenses = await fetchBudgetList({
    year: currentYear,
    month: currentMonth,
  });
  console.log("Fetched expenses:", expenses); // デバッグ用のログ

  return (
    <Suspense fallback={<div>読み込み中...</div>}>
      {expenses.map((expense) => (
        <div key={expense.id} className="expenseItemBox">
          <p>店名:{expense.store_name}</p>
          <span>金額:{expense.amount}</span>
          <span>日付:{expense.date}</span>
        </div>
      ))}
    </Suspense>
  );
}
