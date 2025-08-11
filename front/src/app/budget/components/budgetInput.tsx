"use client";

import { TextBox } from "@/components/elements/textBox/textBox";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { createHeaders } from "@/utils/getCsrf";
import { formattedDate } from "@/utils/formattedDate";

export default function BudgetInput() {
  const router = useRouter();
  const [budgetInput, setBudgetInput] = useState({
    amount: "",
    storeName: "",
    date: "",
    category: "",
    memo: "",
  });
  const [error, setError] = useState("");

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setBudgetInput((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");

    try {
      // 数値型に変換
      const amountNumber = parseInt(budgetInput.amount);

      if (isNaN(amountNumber)) {
        setError("金額は有効な数値である必要があります");
        return;
      }

      const headers = await createHeaders();

      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/expenses`, {
        method: "POST",
        headers: headers,
        body: JSON.stringify({
          amount: amountNumber,
          store_name: budgetInput.storeName,
          date: formattedDate(budgetInput.date),
          category: budgetInput.category,
          memo: budgetInput.memo,
        }),
        credentials: "include",
      });

      if (res.ok) {
        // フォームをリセット
        setBudgetInput({
          amount: "",
          storeName: "",
          date: "",
          category: "",
          memo: "",
        });

        router.push("/budget");
        router.refresh();
        alert("支出を登録しました。");
      } else {
        const errorData = await res.json();
        setError(errorData.error || "支出の登録に失敗しました");
      }
    } catch (err) {
      console.error(err);
      setError("支出の登録中にエラーが発生しました");
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>予算入力</h2>
      {error && <p className="text-red-500 text-sm mb-4">{error}</p>}

      <div>
        <TextBox
          label="金額"
          type="number"
          name="amount"
          value={budgetInput.amount}
          onChange={handleChange}
          placeholder="¥0"
          required
        />

        <TextBox
          label="店名"
          type="text"
          name="storeName"
          value={budgetInput.storeName}
          onChange={handleChange}
          placeholder="店名を入力"
          required
        />

        <TextBox
          label="日付"
          type="date"
          name="date"
          value={budgetInput.date}
          onChange={handleChange}
          required
        />

        <TextBox
          label="カテゴリー"
          type="text"
          name="category"
          value={budgetInput.category}
          onChange={handleChange}
          placeholder="カテゴリーを入力"
          required
        />

        <TextBox
          label="メモ"
          type="text"
          name="memo"
          value={budgetInput.memo}
          onChange={handleChange}
          placeholder="メモを入力（任意）"
        />

        <button type="submit">登録する</button>
      </div>
    </form>
  );
}
