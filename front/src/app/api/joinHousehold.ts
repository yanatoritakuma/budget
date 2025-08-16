import { createHeaders } from "@/utils/getCsrf";

export async function joinHousehold(inviteCode: string): Promise<void> {
  const headers = await createHeaders();
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/household/join`, {
    method: "POST",
    headers: headers,
    body: JSON.stringify({ invite_code: inviteCode }),
    credentials: "include",
  });

  if (!res.ok) {
    const errorData = await res.json();
    throw new Error(errorData.error || "Failed to join household");
  }
}