FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY /cmd/ /app/cmd/
COPY /internal/ /app/internal/
COPY ./init.sql /app/init.sql

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o tendermanagement ./cmd/main.go

FROM alpine:latest

RUN apk update && apk --no-cache add ca-certificates

COPY --from=builder /app/tendermanagement /tendermanagement
COPY --from=builder /app/init.sql /init.sql

EXPOSE 8080

CMD ["/tendermanagement"]
