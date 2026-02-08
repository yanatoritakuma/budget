import { createHeaders } from "@/utils/getCsrf";

type CreateAccountResponse = {
  message: string;
};

export async function lineCreateAccount(): Promise<CreateAccountResponse> {
  const headers = await createHeaders();
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/line/create`, {
    method: 'POST',
    headers: headers,
    credentials: "include",
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || 'Failed to create account');
  }

  return response.json();
}
