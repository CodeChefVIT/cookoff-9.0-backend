FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./dist/app ./cmd/api


FROM alpine:latest 

COPY --from=builder /app/dist/app .

ENV HTTP_PORT=8080

EXPOSE $HTTP_PORT

CMD ["./app"]

