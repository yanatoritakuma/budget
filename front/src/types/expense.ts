export type Expense = {
  id: number;
  user_id: number;
  amount: number;
  store_name: string;
  date: string;
  category: string;
  memo: string;
  created_at: string;
  payer_name?: string;
};