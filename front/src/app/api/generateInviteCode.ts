import { createHeaders } from "@/utils/getCsrf";

export async function generateInviteCode(): Promise<string> {
  const headers = await createHeaders();
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/household/invite-code`, {
    method: "POST",
    headers: headers,
    credentials: "include",
  });

  if (!res.ok) {
    throw new Error("Failed to generate invite code");
  }

  const data = await res.json();
  return data.invite_code;
}