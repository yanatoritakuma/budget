import { Expense } from "@/types/expense";
import { cookies } from "next/headers";

type FetchBudgetListParams = {
  year: number;
  month: number;
  category?: string;
};

export async function fetchBudgetList({
  year,
  month,
  category,
}: FetchBudgetListParams): Promise<Expense[]> {
  const params = new URLSearchParams();
  params.append("year", year.toString());
  params.append("month", month.toString());
  if (category) {
    params.append("category", category);
  }

  const cookieStore = cookies();
  const cookie = cookieStore.toString();

  const res = await fetch(
    `${process.env.NEXT_PUBLIC_API_URL}/expenses?${params.toString()}`,
    {
      method: "GET",
      headers: {
        Cookie: cookie,
      },
      cache: "no-store",
    }
  );

  if (!res.ok) {
    throw new Error("Failed to fetch budget list");
  }

  const resJson: Expense[] = await res.json();

  return resJson;
}
