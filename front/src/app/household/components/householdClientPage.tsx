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
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("招待コードの生成に失敗しました。");
      }
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
      router.push("/budget");
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("世帯への参加に失敗しました。");
      }
    }
  };

  return (
    <div>
      <section className="household-section">
        <h2>現在の世帯メンバー</h2>
        <ul>
          {initialUsers.map((user) => (
            <li key={user.id}>{user.name}</li>
          ))}
        </ul>
      </section>

      <section className="household-section">
        <h2>他のユーザーを招待</h2>
        <button onClick={handleGenerateCode}>招待コードを生成</button>
        {inviteCode && (
          <div className="invite-code-display">
            <p>このコードを他のユーザーと共有してください:</p>
            <code>{inviteCode}</code>
          </div>
        )}
      </section>

      <section className="household-section">
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

      {error && <p className="status-message error">{error}</p>}
      {message && <p className="status-message success">{message}</p>}
    </div>
  );
}
