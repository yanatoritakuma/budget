"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { lineAuthCallback } from "../../api/lineAuthCallback";

export default function LineAuthCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const code = searchParams.get("code");
    const state = searchParams.get("state");

    if (!code || !state) {
      setError("認証コードまたは状態がありません。");
      return;
    }

    const handleLineCallback = async () => {
      try {
        await lineAuthCallback(code, state);

        router.push("/budget");
      } catch (err) {
        setError("LINEログイン中にエラーが発生しました。");
        console.error("LINE Login Callback Error:", err);
      }
    };

    handleLineCallback();
  }, [searchParams, router]);

  if (error) {
    return (
      <div className="line-login-error">
        <p>エラー: {error}</p>
        <button onClick={() => router.push("/")}>ホームに戻る</button>
      </div>
    );
  }

  return null;
}
