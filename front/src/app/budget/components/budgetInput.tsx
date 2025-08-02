"use client";

import { TextBox } from "@/components/elements/textBox/textBox";
import { useState } from "react";

export default function BudgetInput() {
  const [budgetInput, setBudgetInput] = useState({
    amount: "",
    storeName: "",
    date: "",
    category: "",
    memo: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setBudgetInput((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    // TODO: APIにデータを送信する処理を追加
    console.log(budgetInput);
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>予算入力</h2>

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
