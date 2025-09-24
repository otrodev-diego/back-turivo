FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ./cmd/api

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /out/app /app
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app"]

