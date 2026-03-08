.PHONY: run swagger dev-run test

run:
	go run cmd/api/main.go

swagger:
	swag init -g cmd/api/main.go

dev-run: swagger run

COVERAGE_EXCLUDE=mocks|main.go|docs|test

test:
	go test ./... -coverprofile=coverage.tmp
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.tmp coverage.out