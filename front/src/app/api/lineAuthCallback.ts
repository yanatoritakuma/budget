// front/src/app/api/lineAuthCallback.ts

export type LineAuthCallbackResponse = {
  message?: string;
  status: "logged_in" | "unregistered";
  line_name?: string;
  line_picture?: string;
};

type LineAuthCallbackError = {
  error: string;
};

export async function lineAuthCallback(
  code: string,
  state: string,
): Promise<LineAuthCallbackResponse> {
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/line/callback?code=${code}&state=${state}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: 'include',
    },
  );

  if (!response.ok) {
    const errorData: LineAuthCallbackError = await response.json();
    throw new Error(errorData.error || "LINEログインに失敗しました。");
  }

  return response.json();
}
