FROM golang:1.21.5 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/tlsproxy

FROM gcr.io/distroless/cc-debian12

COPY --from=builder /app/tlsproxy /app/tlsproxy
ENTRYPOINT ["/app/tlsproxy"]
