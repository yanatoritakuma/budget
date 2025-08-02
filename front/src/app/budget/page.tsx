import "@/app/budget/styles.scss";
import { fetchLoginUser } from "@/app/api/fetchLoginUser";
import BudgetInput from "@/app/budget/components/budgetInput";
import Link from "next/link";

export default async function Page() {
  const loginUser = await fetchLoginUser();
  return (
    <main className="pageBox">
      <h2>家計簿</h2>
      {loginUser ? (
        <BudgetInput />
      ) : (
        <>
          <p>ログインしてください。</p>
          <Link href="/login">ログインへ</Link>
        </>
      )}
    </main>
  );
}
