"use client";

import { useState } from "react";
import { TLoginUser } from "@/app/api/fetchLoginUser";
import { generateInviteCode } from "@/app/api/generateInviteCode";
import { joinHousehold } from "@/app/api/joinHousehold";
import { useRouter } from "next/navigation";

type Props = {
  initialUsers: TLoginUser[];
};

export default function HouseholdClientPage({ initialUsers }: Props) {
  const router = useRouter();
  const [users, setUsers] = useState<TLoginUser[]>(initialUsers);
  const [inviteCode, setInviteCode] = useState<string>("");
  const [joinCode, setJoinCode] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [message, setMessage] = useState<string>("");

  const handleGenerateCode = async () => {
    setError("");
    setMessage("");
    try {
      const newCode = await generateInviteCode();
      setInviteCode(newCode);
      setMessage("招待コードを生成しました。");
    } catch (err: any) {
      setError(err.message || "招待コードの生成に失敗しました。");
    }
  };

  const handleJoinHousehold = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setMessage("");
    try {
      await joinHousehold(joinCode);
      setMessage("世帯に参加しました！ページを更新しています...");
      // Refresh the page to reflect the new household
      router.refresh();
    } catch (err: any) {
      setError(err.message || "世帯への参加に失敗しました。");
    }
  };

  return (
    <div>
      <section>
        <h2>現在の世帯メンバー</h2>
        <ul>
          {users.map((user) => (
            <li key={user.id}>{user.name}</li>
          ))}
        </ul>
      </section>

      <hr />

      <section>
        <h2>他のユーザーを招待</h2>
        <button onClick={handleGenerateCode}>招待コードを生成</button>
        {inviteCode && (
          <div>
            <p>このコードを他のユーザーと共有してください:</p>
            <code>{inviteCode}</code>
          </div>
        )}
      </section>

      <hr />

      <section>
        <h2>別の世帯に参加</h2>
        <form onSubmit={handleJoinHousehold}>
          <input
            type="text"
            value={joinCode}
            onChange={(e) => setJoinCode(e.target.value)}
            placeholder="招待コードを入力"
            required
          />
          <button type="submit">参加する</button>
        </form>
      </section>

      {error && <p style={{ color: "red" }}>{error}</p>}
      {message && <p style={{ color: "green" }}>{message}</p>}
    </div>
  );
}