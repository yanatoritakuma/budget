"use client";

import { TextBox } from "@/components/elements/textBox/textBox";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { createHeaders } from "@/utils/getCsrf";
import { formattedDate } from "@/utils/formattedDate";
import { TLoginUser } from "@/app/api/fetchLoginUser";
import { SelectBox } from "@/components/elements/selectBox/selectBox";
import { SelectChangeEvent } from "@mui/material/Select";
import "./budgetInput.scss";

type Props = {
  loginUser: TLoginUser;
  householdUsers: TLoginUser[];
};

export default function BudgetInput({ loginUser, householdUsers }: Props) {
  const router = useRouter();
  const [budgetInput, setBudgetInput] = useState({
    amount: "",
    storeName: "",
    date: "",
    category: "",
    memo: "",
    payerId: String(loginUser.id), // Default to the logged-in user
  });
  const [error, setError] = useState("");

  const payerOptions = householdUsers.map((user) => ({
    value: user.id,
    label: user.name,
  }));

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setBudgetInput((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSelectChange = (e: SelectChangeEvent<string | number>) => {
    const { name, value } = e.target;
    setBudgetInput((prev) => ({
      ...prev,
      [name]: String(value),
    }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");

    try {
      const amountNumber = parseInt(budgetInput.amount);
      const payerIdNumber = parseInt(budgetInput.payerId);

      if (isNaN(amountNumber)) {
        setError("金額は有効な数値である必要があります");
        return;
      }
      if (isNaN(payerIdNumber)) {
        setError("有効な支払者が選択されていません");
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
          payer_id: payerIdNumber,
        }),
        credentials: "include",
      });

      if (res.ok) {
        setBudgetInput({
          amount: "",
          storeName: "",
          date: "",
          category: "",
          memo: "",
          payerId: String(loginUser.id),
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
    <form onSubmit={handleSubmit} className="budget-input-form">
      <h2>予算入力</h2>
      {error && <p className="error-message">{error}</p>}

      <div className="form-fields">
        <TextBox
          label="金額"
          type="number"
          name="amount"
          value={budgetInput.amount}
          onChange={handleInputChange}
          placeholder="¥0"
          required
        />

        <TextBox
          label="店名"
          type="text"
          name="storeName"
          value={budgetInput.storeName}
          onChange={handleInputChange}
          placeholder="店名を入力"
          required
        />

        <TextBox
          label="日付"
          type="date"
          name="date"
          value={budgetInput.date}
          onChange={handleInputChange}
          required
        />

        <TextBox
          label="カテゴリー"
          type="text"
          name="category"
          value={budgetInput.category}
          onChange={handleInputChange}
          placeholder="カテゴリーを入力"
          required
        />

        <TextBox
          label="メモ"
          type="text"
          name="memo"
          value={budgetInput.memo}
          onChange={handleInputChange}
          placeholder="メモを入力（任意）"
        />

        <SelectBox
          label="支払った人"
          name="payerId"
          value={budgetInput.payerId}
          onChange={handleSelectChange}
          options={payerOptions}
          required
        />

        <button type="submit" className="submit-button">
          登録する
        </button>
      </div>
    </form>
  );
}
