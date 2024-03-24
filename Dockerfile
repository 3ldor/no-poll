# build stage
FROM golang:1.22-alpine as builder

WORKDIR /app

# download modules as distinct layer
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# build code
COPY *.go ./
RUN go build -ldflags="-s -w" -o no-poll ./...

# runtime stage
FROM alpine

WORKDIR /app
COPY --from=builder /app/no-poll ./

CMD ["./no-poll"]
