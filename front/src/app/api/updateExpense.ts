import { components } from "@/types/api";
import { createHeaders } from "@/utils/getCsrf";

type Expense = components["schemas"]["ExpenseRequest"];
type ExpenseResponse = components["schemas"]["ExpenseResponse"];

export const updateExpense = async (expense: Expense, expenseId: number): Promise<ExpenseResponse> => {
  const headers = await createHeaders();
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/expenses/${expenseId}`, {
    method: 'PUT',
    headers: headers,
    body: JSON.stringify(expense),
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error('Failed to update expense');
  }

  return response.json();
};
