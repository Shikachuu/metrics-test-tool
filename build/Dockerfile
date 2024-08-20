FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/ cmd/
COPY magefile.go magefile.go
RUN go run github.com/magefile/mage build:binary

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app/bin/ /app/bin/
ENTRYPOINT ["/app/bin/metrics-tool"]
