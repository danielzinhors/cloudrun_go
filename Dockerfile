FROM golang:1.21.6 as builder
WORKDIR /server
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o temp_server cmd/server/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /server/temp_server .
ENTRYPOINT ["./temp_server"]