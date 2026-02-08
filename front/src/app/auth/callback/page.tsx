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
        const data = await lineAuthCallback(code, state);

        if (data.status === "unregistered") {
          const params = new URLSearchParams();
          if (data.line_name) params.set("name", data.line_name);
          if (data.line_picture) params.set("picture", data.line_picture);
          router.push(`/auth/line-link?${params.toString()}`);
        } else {
          router.push("/budget");
        }
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
