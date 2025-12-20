import { createHeaders } from "@/utils/getCsrf";

export const deleteExpense = async (expenseId: number): Promise<void> => {
  const headers = await createHeaders();
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/expenses/${expenseId}`, {
    method: 'DELETE',
    headers: headers,
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error('Failed to delete expense');
  }

  return;
};
