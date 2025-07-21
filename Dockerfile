FROM golang:1.24-bullseye AS builder

WORKDIR /app
COPY . .

RUN apt-get update && apt-get install -y \
    make \
    gcc \
    librdkafka-dev \
    pkg-config

RUN if [ ! -f "go.mod" ]; then \
        make init; \
    fi

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o /app/ms_exchange ./cmd/ms_exchange

FROM debian:bullseye-slim AS runner

WORKDIR /app
COPY --from=builder /app/ms_exchange /app/

RUN apt-get update && apt-get install -y \
    librdkafka++1 && \
    rm -rf /var/lib/apt/lists/*

CMD ["/app/ms_exchange"]