import "@/app/auth/styles.scss";
import { fetchLoginUser } from "@/app/api/fetchLoginUser";
import Auth from "@/app/auth/components/auth";
import Link from "next/link";
import { createHeaders } from "@/utils/getCsrf";

export default async function Page() {
  const loginUser = await fetchLoginUser();

  let lineAuthUrl: string | null = null;
  // Fetch auth_url only if not logged in
  if (!loginUser?.id) {
    try {
      const headers = await createHeaders();
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/line/login`,
        {
          method: "GET",
          headers: headers,
          cache: "no-store",
          credentials: "include",
        },
      );

      if (res.ok) {
        const data = await res.json();
        if (data.auth_url) {
          lineAuthUrl = data.auth_url;
        }
      } else {
        console.error(
          "Failed to fetch LINE auth URL:",
          res.status,
          res.statusText,
        );
      }
    } catch (err) {
      console.error("Error fetching LINE auth URL:", err);
    }
  }

  return (
    <main className="pageBox">
      <div className="">
        {loginUser?.id !== undefined ? (
          <div className="auth__loggedIn">
            <h3>ログイン済みです。</h3>
            <Link prefetch={false} href="/">
              ホームへ
            </Link>
          </div>
        ) : (
          <Auth lineAuthUrl={lineAuthUrl} />
        )}
      </div>
    </main>
  );
}
