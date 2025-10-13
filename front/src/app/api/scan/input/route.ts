import { NextRequest, NextResponse } from "next/server";

export async function POST(req: NextRequest) {
  try {
    const { image } = await req.json();

    // data URL → Blob
    const fetchRes = await fetch(image);
    const blob = await fetchRes.blob();

    const formData = new FormData();
    const filename = "receipt." + blob.type.split("/")[1];
    formData.append("file", blob, filename);

    // 画像アップロード
    const uploadRes = await fetch("https://api.tabscanner.com/api/2/process", {
      method: "POST",
      headers: {
        apikey: process.env.TABSCANNER_API_KEY || "",
      },
      body: formData,
    });

    const uploadData = await uploadRes.json();

    if (!uploadRes.ok || !uploadData.success) {
      return NextResponse.json(
        { error: uploadData.message || "Upload failed" },
        { status: uploadRes.status }
      );
    }

    // すぐにトークンを含むデータを返す
    return NextResponse.json(uploadData);
  } catch (error) {
    console.error("Error in /api/scan/input:", error);
    return NextResponse.json(
      { error: "Internal Server Error" },
      { status: 500 }
    );
  }
}
