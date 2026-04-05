.PHONY: run swagger dev-run test

run:
	go run cmd/api/main.go

swagger:
	swag init -g cmd/api/main.go

dev-run: swagger run

COVERAGE_EXCLUDE=mocks|main.go|docs|test|cmd|pkg|internal/api
COVERAGE_THRESHOLD=80

test:
	go test ./... -coverprofile=coverage.tmp
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	
	@{ \
	total=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$total%"; \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		rm coverage.tmp coverage.out; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi; \
	}

	# Generate HTML report
	go tool cover -html=coverage.out -o coverage.html

	# Cleanup
	rm coverage.tmp coverage.out

up:
	docker compose down
	docker compose up --build -d

down:
	docker compose down

docker-build:
	docker build -t bookmark_service:dev .