generate-go:
	oapi-codegen -generate types -o internal/api/types.gen.go -package api openapi.yaml
	oapi-codegen -generate chi-server -o internal/api/server.gen.go -package api openapi.yaml

build-lambda:
	# AWS Lambda向けのGoバイナリをビルドします。
	# glibcのバージョン問題に対応するため、linux/amd64向けにクロスコンパイルし、静的リンクします。
	# 出力ファイル名はAWS Lambdaの規約に従い 'bootstrap' とします。
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap -ldflags="-s -w" .
	# 既存のzipファイルを削除し、クリーンな状態で再作成します。
	rm -f lambda-handler.zip
	# ビルドしたバイナリをzipファイルに圧縮します。
	zip lambda-handler.zip bootstrap
	# 一時的な 'bootstrap' ファイルを削除します。
	rm bootstrap

