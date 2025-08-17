import { format, parseISO } from "date-fns";

/**
 * Parses an ISO date string and formats it for display.
 * @param date - The ISO date string (e.g., "2025-08-16T10:00:00Z").
 * @param formatString - The desired output format (defaults to "yyyy/MM/dd").
 * @returns The formatted date string.
 */
export const formatDateForDisplay = (date: string, formatString = "yyyy/MM/dd") => {
  try {
    const parsedDate = parseISO(date);
    return format(parsedDate, formatString);
  } catch (error) {
    console.error("Invalid date format for display:", date, error);
    // Return the original string or a sensible default if parsing fails
    return date;
  }
};