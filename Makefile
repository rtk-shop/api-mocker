go:
	@go run ./cmd/main.go -env=local
gen:
	@go tool gqlgenc
oapi:
	@go tool oapi-codegen -config oapi-codegen.yaml api.yaml
schema:
	@get-graphql-schema http://localhost:8080/graphql > schema.graphql
