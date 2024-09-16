FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY /cmd/ /app/cmd/
COPY /internal/ /app/internal/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o tendermanagement ./cmd/main.go

FROM alpine:latest

RUN apk update && apk --no-cache add ca-certificates

COPY --from=builder /app/tendermanagement /tendermanagement

ENV SERVER_ADDRESS=0.0.0.0:8080
ENV POSTGRES_USERNAME=cnrprod1725727312-team-79255
ENV POSTGRES_PASSWORD=cnrprod1725727312-team-79255
ENV POSTGRES_HOST=rc1b-5xmqy6bq501kls4m.mdb.yandexcloud.net
ENV POSTGRES_PORT=6432
ENV POSTGRES_DATABASE=cnrprod1725727312-team-79255

EXPOSE 8080

CMD ["/tendermanagement"]