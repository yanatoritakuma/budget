package utils

import "net"

// extractHostname は URL から ホスト名のみを抽出します。
// Cookie の domain パラメータに渡すために使用します。
func ExtractHostname(urlStr string) string {
	if urlStr == "" {
		return "" // domain 未設定（現在のホストに限定）
	}

	// URL をパースしてホスト部分を抽出
	// "http://localhost:3000" → "localhost"
	// "https://example.com" → "example.com"
	host, _, err := net.SplitHostPort(urlStr)
	if err != nil {
		// ポート番号がない場合、そのままホスト名
		// スキームを削除
		if len(urlStr) > 7 && (urlStr[:7] == "http://" || urlStr[:8] == "https://") {
			if urlStr[:7] == "http://" {
				return urlStr[7:]
			}
			return urlStr[8:]
		}
		return urlStr
	}
	return host
}
