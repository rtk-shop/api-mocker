go:
	@go run ./cmd/main.go -env=local
gen:
	@go tool gqlgenc
schema:
	@get-graphql-schema http://localhost:8080/graphql > schema.graphql
