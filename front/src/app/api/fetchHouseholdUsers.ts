import { headers } from "next/headers";
import { TLoginUser } from "./fetchLoginUser";

export async function fetchHouseholdUsers(): Promise<TLoginUser[]> {
  const headersList = await headers();
  const token = headersList.get("cookie");

  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/household/users`, {
    method: "GET",
    headers: {
      ...headers,
      Cookie: `${token}`,
    },
    cache: "no-store",
    credentials: "include",
  });

  if (!res.ok) {
    // This will activate the closest `error.js` Error Boundary
    throw new Error('Failed to fetch household users');
  }

  const users: TLoginUser[] = await res.json();
  return users;
}