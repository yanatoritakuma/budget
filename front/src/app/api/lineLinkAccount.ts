import { createHeaders } from "@/utils/getCsrf";

type LinkAccountResponse = {
  message: string;
};

export async function lineLinkAccount(email: string, password: string): Promise<LinkAccountResponse> {
  const headers = await createHeaders();
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/line/link`, {
    method: 'POST',
    headers: headers,
    body: JSON.stringify({ email, password }),
    credentials: "include",
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || 'Failed to link account');
  }

  return response.json();
}
