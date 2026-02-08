# **予算超過通知機能 設計書**

## **概要**

ユーザーが支出を登録した際、その世帯の当月の合計支出が設定された予算を超過した場合に、LINEで通知を送信する。 通知は月ごとに1回のみ行い、重複した通知を防ぐ。

## **アーキテクチャ**

1. 支出登録（`POST /expenses`）
2. 予算チェックロジック（Usecase）
3. 条件合致時、SQSへイベント送信
4. （後段）通知LambdaがSQSをトリガーに起動
    - `household_id` を元に、DBからその世帯に所属する全ユーザーの `LineUserID` を取得。
    - LINE Messaging API を使用して通知を送信。

## **データモデル**

### **Budget テーブル**

月ごとの予算と通知状態を管理する。

| カラム名 | 型 | 説明 |
| --- | --- | --- |
| id | SERIAL | プライマリキー |
| household_id | INT | 世帯ID（外部キー） |
| year_month | VARCHAR(7) | 対象年月（"YYYY/MM" 形式） |
| amount | INT | 予算額 |
| notified_at | TIMESTAMP | 通知送信日時（NULLの場合は未通知） |
| created_at | TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | 更新日時 |

## **通知判定ロジック**

支出登録時（`CreateExpense`）に以下のフローで判定を行う。

1. 支出の登録年月（`YYYY/MM`）を取得。
2. `household_id` と `year_month` をキーに `budgets` テーブルを検索。
3. 予算レコードが存在し、かつ `notified_at` が NULL である場合のみ続行。
4. その世帯の当月の合計支出額（今回分を含む）を計算。
5. `合計支出額 > 予算額` の場合：
    - SQSに通知イベントを送信。
    - `budgets` テーブルの `notified_at` を現在時刻で更新。

## **SQS メッセージ構造**

通知Lambdaが必要とする情報を JSON 形式で送信する。

```json
{
  "type": "BUDGET_EXCEEDED",
  "household_id": 1,
  "year_month": "2026/02",
  "budget_amount": 30000,
  "current_amount": 30500
}
```