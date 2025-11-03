generate-go:
	oapi-codegen -generate types -o internal/api/types.gen.go -package api openapi.yaml
	oapi-codegen -generate chi-server -o internal/api/server.gen.go -package api openapi.yaml
