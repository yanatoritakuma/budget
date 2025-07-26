import "@/app/style.scss";
import Link from "next/link";

export default function Home() {
  return (
    <div>
      <main className="pageBox">
        <h1>家計簿</h1>
        <Link href="/auth">ログイン</Link>
      </main>
    </div>
  );
}
