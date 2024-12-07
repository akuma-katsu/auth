FROM golang:1.22.5-alpine as builder

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . ./
RUN go build -o ./auth  backend/cmd/main.go

FROM alpine:latest

COPY --from=builder /app/auth /
COPY  .env /

EXPOSE 8080

CMD ["/auth"]