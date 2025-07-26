import "@/app/auth/styles.scss";
import { fetchLoginUser } from "@/app/api/fetchLoginUser";
import Auth from "@/app/auth/components/auth";
import Link from "next/link";

export default async function Page() {
  const loginUser = await fetchLoginUser();
  return (
    <main className="pageBox">
      <div className="">
        {loginUser.id !== undefined ? (
          <div className="auth__loggedIn">
            <h3>ログイン済みです。</h3>
            <Link prefetch={false} href="/">
              ホームへ
            </Link>
          </div>
        ) : (
          <Auth />
        )}
      </div>
    </main>
  );
}
