import "@/app/budget/styles.scss";
import { fetchLoginUser } from "@/app/api/fetchLoginUser";
import BudgetInput from "@/app/budget/components/budgetInput/budgetInput";
import Link from "next/link";
import BudgetList from "@/app/budget/components/budgetList/budgetList";

export default async function Page() {
  const loginUser = await fetchLoginUser();
  return (
    <main className="pageBox">
      <h2>家計簿</h2>
      {loginUser ? (
        <>
          <BudgetInput />
          <BudgetList />
        </>
      ) : (
        <>
          <p>ログインしてください。</p>
          <Link href="/login">ログインへ</Link>
        </>
      )}
    </main>
  );
}
