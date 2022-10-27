FROM golang:1.19.2-alpine as builder
WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 go build -o ./tunnel/tunnel.bin ./tunnel/

FROM scratch
COPY --from=builder /build/tunnel/tunnel.bin /tunnel.bin
ENTRYPOINT ["/tunnel.bin"]