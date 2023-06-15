FROM golang:1.19-buster as builder

ENV OUTPUT_PATH=/tmp/github-metrics

WORKDIR /app

COPY . .

RUN make install-deps

RUN go build -v -o github-metrics

FROM debian:buster-slim

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY config.yaml /app/config.yaml
COPY --from=builder /app/github-metrics /app/github-metrics

VOLUME /tmp/github-metrics

ENTRYPOINT ["sh", "-c", "/app/github-metrics -out ${OUTPUT_PATH}" ]

