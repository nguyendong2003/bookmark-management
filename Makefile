IMG_NAME=dongcoi14122003/bookmark_service
GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null)
BRANCH  := $(shell git rev-parse --abbrev-ref HEAD)

# Default tag (fallback)
IMG_TAG := latest

# If branch is main, set tag to dev
ifeq ($(BRANCH),main)
    IMG_TAG := dev
endif

# If git tag is set, use it
ifneq ($(GIT_TAG),)
    IMG_TAG := $(GIT_TAG)
endif

# Export the tag to be used in other targets
export IMG_TAG

.PHONY: dev-run swagger run-app test docker-test docker-build build

dev-run:
	swag init -g cmd/api/main.go
	go run cmd/api/main.go

swagger:
	swag init -g cmd/api/main.go

run-app:
	go run cmd/api/main.go

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

COVERAGE_EXCLUDE=mocks|main.go|docs|test|cmd|pkg|internal/api
COVERAGE_THRESHOLD = 80

test:
	go test ./... -coverprofile=coverage.tmp -coverpkg=./... -covermode=atomic -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		rm -f coverage.tmp coverage.out; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
		rm -f coverage.tmp coverage.out; \
	fi

COVERAGE_FOLDER=./coverage

docker-test:
	mkdir -p $(COVERAGE_FOLDER)
	docker buildx build --build-arg COVERAGE_EXCLUDE="$(COVERAGE_EXCLUDE)" --target test -t bookmark-service-test:dev --output $(COVERAGE_FOLDER) .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

docker-build:
	docker build -t $(IMG_NAME):$(IMG_TAG) .

docker-release:
	docker push $(IMG_NAME):$(IMG_TAG)

DOCKER_USERNAME ?=
DOCKER_PASSWORD ?=

docker-login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin