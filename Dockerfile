FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o viandas-api ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/viandas-api .
COPY --from=builder /app/settings/settings.yml ./settings/settings.yml

CMD ["./viandas-api"]
