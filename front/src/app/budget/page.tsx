import "@/app/budget/styles.scss";
import { fetchLoginUser } from "@/app/api/fetchLoginUser";
import BudgetInput from "@/app/budget/components/budgetInput/budgetInput";
import Link from "next/link";
import BudgetList from "@/app/budget/components/budgetList/budgetList";
import { fetchHouseholdUsers } from "@/app/api/fetchHouseholdUsers";

export default async function Page() {
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

  return (
    <main className="pageBox">
      <h2>家計簿</h2>
      <BudgetInput loginUser={loginUser} householdUsers={householdUsers} />
      <BudgetList />
    </main>
  );
}
