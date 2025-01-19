FROM golang:1.23 AS compiler
WORKDIR /app
COPY . .
RUN go build -a -tags musl --ldflags "-linkmode external -extldflags '-static' -s -w" -o roundguard ./cmd/
FROM scratch
COPY --from=compiler /app/roundguard .

CMD ["./roundguard", "echo", "start"]
