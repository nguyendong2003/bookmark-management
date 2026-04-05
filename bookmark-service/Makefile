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
	
	# Tạm thời vô hiệu hóa kiểm tra ngưỡng coverage
	# if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	# 	echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	# 	exit 1; \
	# else \
	# 	echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	# fi

	# Xóa file coverage.tmp và coverage.out
	rm coverage.tmp coverage.out

up:
	docker compose down
	docker compose up --build -d

down:
	docker compose down

docker-build:
	docker build -t bookmark_service:dev .