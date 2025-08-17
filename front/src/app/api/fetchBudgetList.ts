import { headers } from "next/headers";

type FetchBudgetListParams = {
  year: number;
  month: number;
  category?: string;
};

type ExpenseResponse = {
  amount: number;
  category: Date;
  created_at: Date;
  date: string;
  id: number;
  memo: string;
  store_name: string;
  user_id: number;
  payer_name?: string; // Add optional payer_name
};

export async function fetchBudgetList({
  year,
  month,
  category,
}: FetchBudgetListParams) {
  const headersList = await headers();
  const token = headersList.get("cookie");
  const params = new URLSearchParams();
  params.append("year", year.toString());
  params.append("month", month.toString());
  if (category) {
    params.append("category", category);
  }
  const res = await fetch(
    `${process.env.NEXT_PUBLIC_API_URL}/expenses?${params.toString()}`,
    {
      method: "GET",
      headers: {
        ...headers,
        Cookie: `${token}`,
      },
      cache: "no-store",
      // cache: "force-cache",
      credentials: "include",
    }
  );

  const resJson: ExpenseResponse[] = await res.json();

  return resJson;
}
