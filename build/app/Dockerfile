FROM golang:1.23.6-alpine3.20 as builder

WORKDIR /app

RUN apk update && \
    apk add --no-cache git

ADD . .

RUN go mod download

RUN go build -o main ./app/cmd/http/main.go

FROM alpine:3.20 AS running

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata libc6-compat
ENV TZ=Asia/Bangkok

ENV APPENV=local

COPY --from=builder /app/main .
COPY /config/local.config.yaml ./config/local.config.yaml
COPY /config/.local.env ./config/.local.env

RUN apk add --no-cache libcap && \
    setcap 'cap_net_bind_service=+ep' /app/main && \
    apk del libcap


EXPOSE 8080

ENTRYPOINT ["/app/main"]
