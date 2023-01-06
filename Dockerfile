
FROM golang:1.18 as builder

WORKDIR /app
COPY go.mod /app
COPY src/main /app

RUN go get -d -v

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o app .

FROM gcr.io/distroless/static

COPY --from=builder /app/app /app

ENTRYPOINT ["/app"]