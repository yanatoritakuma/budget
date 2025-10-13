import { NextRequest, NextResponse } from "next/server";

async function fetchResult(token: string) {
  const resultRes = await fetch(
    `https://api.tabscanner.com/api/result/${token}`,
    {
      method: "GET",
      headers: {
        apikey: process.env.TABSCANNER_API_KEY || "",
      },
      cache: "no-store", // Important for polling
    }
  );
  if (!resultRes.ok) {
    throw new Error(`Failed to fetch result: ${resultRes.statusText}`);
  }
  return resultRes.json();
}

export async function GET(req: NextRequest) {
  const { searchParams } = new URL(req.url);
  const token = searchParams.get("token");

  if (!token) {
    return NextResponse.json({ error: "Token is required" }, { status: 400 });
  }

  const MAX_ATTEMPTS = 15;
  const DELAY_MS = 2000; // 2 seconds

  for (let i = 0; i < MAX_ATTEMPTS; i++) {
    try {
      const resultData = await fetchResult(token);

      if (resultData.status === "done") {
        return NextResponse.json(resultData);
      }
      // If status is 'pending' or something else, wait and loop again

      await new Promise((resolve) => setTimeout(resolve, DELAY_MS));
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (error: any) {
      return NextResponse.json(
        { error: "Failed to poll for results", details: error.message },
        { status: 500 }
      );
    }
  }

  return NextResponse.json({ error: "Polling timed out" }, { status: 504 }); // 504 Gateway Timeout
}
