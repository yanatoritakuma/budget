"use client";

import { useState } from "react"; // Removed useRef
import "@/app/auth/components/auth.scss";
import { createHeaders } from "@/utils/getCsrf";
import { TextBox } from "@/components/elements/textBox/textBox";
import { ButtonBox } from "@/components/elements/buttonBox/buttonBox";
import { useRouter } from "next/navigation";

type AuthProps = {
  lineAuthUrl: string | null;
};

export default function Auth({ lineAuthUrl }: AuthProps) {
  console.log("lineAuthUrl:", lineAuthUrl);
  const router = useRouter();
  const [authState, setAuthState] = useState({
    mail: "",
    password: "",
    name: "",
  });
  const [isLogin, setIsLogin] = useState(true);

  const email = authState.mail;
  const password = authState.password;
  const name = authState.name;

  const onClickAuth = async () => {
    const headers = await createHeaders();
    if (isLogin) {
      try {
        // await new Promise((resolve) => setTimeout(resolve, 3000));
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/login`, {
          method: "POST",
          headers: headers,
          body: JSON.stringify({ email, password }),
          cache: "no-store",
          credentials: "include",
        });

        if (res.ok) {
          router.push("/");
          router.refresh();
          alert("ログインしました。");
        } else {
          alert("ログイン失敗しました。");
        }
      } catch (err) {
        console.error(err);
      }
    } else {
      try {
        // await new Promise((resolve) => setTimeout(resolve, 3000));
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/signup`, {
          method: "POST",
          headers: headers,
          body: JSON.stringify({ email, password, name }),
          cache: "no-store",
          credentials: "include",
        });

        if (res.ok) {
          const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/login`, {
            method: "POST",
            headers: headers,
            body: JSON.stringify({ email, password }),
            cache: "no-store",
            credentials: "include",
          });

          if (res.ok) {
            router.push("/");
            router.refresh();

            alert("アカウント作成しました。");
          } else {
            alert("アカウント失敗しました。");
          }
        }
      } catch (err) {
        console.error(err);
      }
    }
  };

  // Removed handleLineLogin function entirely

  return (
    <section className="authInputBox">
      <h2>{isLogin ? "ログイン" : "新規登録"}</h2>
      <div className="authInputBox__inputBox">
        <TextBox
          label="メールアドレス"
          value={authState.mail}
          onChange={(e) =>
            setAuthState({
              ...authState,
              mail: e.target.value,
            })
          }
          className="authInputBox__input"
          fullWidth
        />
        <TextBox
          label="パスワード"
          value={authState.password}
          onChange={(e) =>
            setAuthState({
              ...authState,
              password: e.target.value,
            })
          }
          password
          className="authInputBox__input"
          fullWidth
        />
        {!isLogin && (
          <TextBox
            label="名前"
            value={authState.name}
            onChange={(e) =>
              setAuthState({
                ...authState,
                name: e.target.value,
              })
            }
            className="authInputBox__input"
            fullWidth
          />
        )}

        <span
          className="authInputBox__text"
          onClick={() => setIsLogin(!isLogin)}
        >
          {isLogin
            ? "アカウントをまだ作成ではない方はこちら"
            : "アカウントをお持ちの方はこちら"}
        </span>
        <ButtonBox onClick={() => onClickAuth()}>
          {isLogin ? "ログイン" : "登録"}
        </ButtonBox>

        {lineAuthUrl ? (
          <a href={lineAuthUrl} target="_self" rel="noopener noreferrer">
            <ButtonBox>LINEでログイン</ButtonBox>
          </a>
        ) : (
          <ButtonBox disabled>LINEでログイン (URL取得中...)</ButtonBox>
        )}
      </div>
    </section>
  );
}
