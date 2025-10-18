import "@/app/budget/styles.scss";
import { fetchLoginUser } from "@/app/api/fetchLoginUser";
import BudgetInput from "@/app/budget/components/budgetInput/budgetInput";
import Link from "next/link";
import BudgetList from "@/app/budget/components/budgetList/budgetList";
import { fetchHouseholdUsers } from "@/app/api/fetchHouseholdUsers";
import { fetchBudgetList } from "@/app/api/fetchBudgetList";
import { Expense } from "@/types/expense";

type PageProps = {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  searchParams?: any;
};

export default async function Page({ searchParams }: PageProps) {
  const loginUser = await fetchLoginUser();

  if (!loginUser) {
    return (
      <main className="pageBox">
        <h2>家計簿</h2>
        <div className="login-prompt">
          <p>ログインしてください。</p>
          <Link href="/auth" className="login-link">
            ログインへ
          </Link>
        </div>
      </main>
    );
  }

  // Only fetch household users if the user is logged in
  const householdUsers = await fetchHouseholdUsers();

  const currentDate = new Date();
  const year = searchParams?.year
    ? parseInt(searchParams.year)
    : currentDate.getFullYear();
  const month = searchParams?.month
    ? parseInt(searchParams.month)
    : currentDate.getMonth() + 1;

  const expenses: Expense[] | null = await fetchBudgetList({ year, month });

  return (
    <main className="pageBox">
      <h2>家計簿</h2>
      <div className="page-nav">
        <Link href="/household">世帯管理へ</Link>
      </div>
      <BudgetInput loginUser={loginUser} householdUsers={householdUsers} />
      <BudgetList expenses={expenses} />
    </main>
  );
}
