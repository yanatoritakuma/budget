import { format } from "date-fns";

export const getCurrentDateFormatted = (): string => {
  return format(new Date(), "yyyy-MM-dd");
};
