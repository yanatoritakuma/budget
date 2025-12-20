import React, { useState, useEffect } from "react";
import "./editModal.scss";
import { components } from "@/types/api";

type Expense = components["schemas"]["ExpenseResponse"];
type User = components["schemas"]["UserResponse"];

interface EditModalProps {
  onClose: () => void;
  expense: Expense | null;
  users: User[];
  onSave: (updatedExpense: Expense) => void;
}

const EditModal: React.FC<EditModalProps> = ({
  onClose,
  expense,
  users,
  onSave,
}) => {
  const [amount, setAmount] = useState(0);
  const [storeName, setStoreName] = useState("");
  const [transactionDate, setTransactionDate] = useState("");
  const [category, setCategory] = useState("");
  const [memo, setMemo] = useState("");
  const [payerId, setPayerId] = useState("");

  useEffect(() => {
    if (expense) {
      setAmount(expense.amount);
      setStoreName(expense.store_name);
      setTransactionDate(new Date(expense.date).toISOString().split("T")[0]);
      setCategory(expense.category);
      setMemo(expense.memo || "");
      setPayerId(expense.user_id.toString());
    }
  }, [expense]);

  const handleSave = () => {
    if (expense) {
      const updatedExpense = {
        ...expense,
        amount,
        store_name: storeName,
        date: `${transactionDate}T00:00:00Z`,
        category,
        memo,
        user_id: parseInt(payerId, 10),
      };
      onSave(updatedExpense);
    }
  };

  return (
    <div className="editModal">
      <div className="modalContent">
        <h2>支出を編集</h2>
        <div className="form">
          <div className="formRow">
            <label className="label">金額</label>
            <input
              type="number"
              className="input"
              value={amount}
              onChange={(e) => setAmount(Number(e.target.value))}
            />
          </div>
          <div className="formRow">
            <label className="label">店名</label>
            <input
              type="text"
              className="input"
              value={storeName}
              onChange={(e) => setStoreName(e.target.value)}
            />
          </div>
          <div className="formRow">
            <label className="label">日付</label>
            <input
              type="date"
              className="input"
              value={transactionDate}
              onChange={(e) => setTransactionDate(e.target.value)}
            />
          </div>
          <div className="formRow">
            <label className="label">カテゴリー</label>
            <input
              type="text"
              className="input"
              value={category}
              onChange={(e) => setCategory(e.target.value)}
            />
          </div>
          <div className="formRow">
            <label className="label">メモ</label>
            <input
              type="text"
              className="input"
              value={memo}
              onChange={(e) => setMemo(e.target.value)}
            />
          </div>
          <div className="formRow">
            <label className="label">支払い者</label>
            <select
              className="input"
              value={payerId}
              onChange={(e) => setPayerId(e.target.value)}
            >
              {users.map((user) => (
                <option key={user.id} value={user.id}>
                  {user.name}
                </option>
              ))}
            </select>
          </div>
        </div>
        <div className="buttonContainer">
          <button className="button cancelButton" onClick={onClose}>
            キャンセル
          </button>
          <button className="button saveButton" onClick={handleSave}>
            保存
          </button>
        </div>
      </div>
    </div>
  );
};


export default EditModal;
