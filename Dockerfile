FROM golang:1.19-alpine3.18 as builder

RUN apk update && apk add --no-cache git

ENV OUTPUT_PATH=/tmp/github-metrics

WORKDIR $GOPATH/src/pkg/app/

COPY . .

RUN go get -d -v

RUN go build -v -o /go/bin/github-metrics

FROM alpine:3.18

RUN apk update && apk add --no-cache git

COPY config.yaml /app/config.yaml
COPY --from=builder /go/bin/github-metrics /app/github-metrics

VOLUME /tmp/github-metrics

ENTRYPOINT ["sh", "-c", "/app/github-metrics -out ${OUTPUT_PATH}" ]

