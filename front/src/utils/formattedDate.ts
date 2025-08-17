import { format, parseISO } from "date-fns";

// date-fnsを使用して日付を適切なフォーマットに変換
export const formattedDate = (date: string) => {
  return format(parseISO(date), "yyyy-MM-dd'T'HH:mm:ss'Z'");
};
