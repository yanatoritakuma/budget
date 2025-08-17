export const createHeaders = async () => {
  const headers = new Headers();
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/csrf`, {
    credentials: "include",
  });
  const data = await response.json();
  headers.append("Content-Type", "application/json");
  headers.append("X-CSRF-Token", data.csrf_token);
  // document.cookie = `_csrf=${data.csrf_token}`;
  headers.append("Cookie", `_csrf=${data.csrf_token}`);
  headers.append("Accept", "application/json, text/plain, */*");
  return headers;
};
