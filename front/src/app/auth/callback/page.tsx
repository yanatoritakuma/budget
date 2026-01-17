"use client";

import { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { createHeaders } from "@/utils/getCsrf";

export default function LineLoginCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const code = searchParams.get("code");
    const state = searchParams.get("state");

    if (!code || !state) {
      console.error("LINEログインコールバック: codeまたはstateが見つかりません。");
      alert("LINEログインに失敗しました。必要な情報が不足しています。");
      router.push("/auth"); // ログインページに戻す
      return;
    }

    const handleLineLoginCallback = async (authCode: string, authState: string) => {
      try {
        const headers = await createHeaders(); // CSRFトークンなどを含むヘッダーを生成
        const res = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/line/callback?code=${authCode}&state=${authState}`,
          {
            method: "GET", // バックエンドAPIはGETリクエストを受け付ける
            headers: headers,
            credentials: "include", // Cookieの送受信を有効にする
          }
        );

        if (res.ok) {
          // バックエンドがリダイレクトを返すため、ここではレスポンスボディを解析しない
          // リダイレクトはブラウザによって自動的に処理される
          // router.push("/budget"); // バックエンドからのリダイレクトを待つ
        } else {
          const errorData = await res.json();
          console.error("LINEログインコールバック失敗:", errorData);
          alert(`LINEログインに失敗しました: ${errorData.error || "不明なエラー"}`);
          router.push("/auth"); // ログインページに戻す
        }
      } catch (err) {
        console.error("LINEログインコールバック中にエラーが発生しました:", err);
        alert("LINEログイン中に予期せぬエラーが発生しました。");
        router.push("/auth"); // ログインページに戻す
      }
    };

    handleLineLoginCallback(code, state);
  }, [searchParams, router]);

  return null;
}
