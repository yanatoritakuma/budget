"use client";

import { useState, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { lineLinkAccount } from "../../api/lineLinkAccount";
import { lineCreateAccount } from "../../api/lineCreateAccount";
import styles from "./lineLink.module.scss";

function LineLinkContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const name = searchParams.get("name");
  const picture = searchParams.get("picture");

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleLink = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      await lineLinkAccount(email, password);
      router.push("/budget");
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("アカウント連携に失敗しました。");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async () => {
    setError(null);
    setLoading(true);

    try {
      await lineCreateAccount();
      router.push("/budget");
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("アカウント作成に失敗しました。");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.card}>
        <div className={styles.profile}>
          {/* eslint-disable-next-line @next/next/no-img-element */}
          {picture && <img src={picture} alt="LINE Profile" />}
          <h2>{name ? `${name}さんとして認証中` : "LINE認証済み"}</h2>
          <p>アカウントを連携するか、新規登録してください。</p>
        </div>

        {error && <p className={styles.error}>{error}</p>}

        <div className={styles.options}>
          <form className={styles.form} onSubmit={handleLink}>
            <h3>既存のアカウントと連携</h3>
            <input
              type="email"
              placeholder="メールアドレス"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
            <input
              type="password"
              placeholder="パスワード"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            <button
              type="submit"
              className={`${styles.button} ${styles.primary}`}
              disabled={loading}
            >
              {loading ? "処理中..." : "連携してログイン"}
            </button>
          </form>

          <div className={styles.divider}>
            <span>または</span>
          </div>

          <div className={styles.form}>
            <h3>新規アカウント作成</h3>
            <button
              type="button"
              className={`${styles.button} ${styles.secondary}`}
              onClick={handleCreate}
              disabled={loading}
            >
              {loading ? "処理中..." : "新規登録してログイン"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function LineLinkPage() {
  return (
    <Suspense fallback={<div className={styles.container}>読み込み中...</div>}>
      <LineLinkContent />
    </Suspense>
  );
}
