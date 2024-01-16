FROM golang:1.19.2-alpine as builder
WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 go build -o ./webfleet/tunnel.bin ./webfleet/

FROM scratch
COPY --from=builder /build/webfleet/tunnel.bin /tunnel.bin
ENTRYPOINT ["/tunnel.bin"]