# Base stage
FROM golang:alpine AS base

RUN mkdir -p /opt/app

WORKDIR /opt/app

RUN apk add build-base

COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go mod download

COPY . .

# Build stage
FROM base AS build

RUN GOOS=linux go build -tags musl -ldflags "-w -s" -o bookmark_service cmd/api/main.go

# Test execution stage
FROM base AS test-exec

ARG _outputdir="/tmp/coverage"
ARG COVERAGE_EXCLUDE

RUN mkdir -p ${_outputdir} && \
    go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1 && \
    grep -v -E "${COVERAGE_EXCLUDE}" coverage.tmp > ${_outputdir}/coverage.out && \
    go tool cover -html=${_outputdir}/coverage.out -o ${_outputdir}/coverage.html

# Store test coverage output
FROM scratch AS test
ARG _outputdir="/tmp/coverage"
COPY --from=test-exec ${_outputdir}/coverage.out /
COPY --from=test-exec ${_outputdir}/coverage.html /

# Final stage
FROM alpine:3.20 AS final

ENV TZ=Asia/Ho_Chi_Minh
WORKDIR /app

# Create user non-root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary + docs
COPY --from=build /opt/app/bookmark_service /app/bookmark_service
COPY --from=build /opt/app/docs /app/docs

# Set ownership
RUN chown -R appuser:appgroup /app

# Set timezone
RUN apk add --no-cache tzdata && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

# Run as non-root user
USER appuser

CMD ["/app/bookmark_service"]