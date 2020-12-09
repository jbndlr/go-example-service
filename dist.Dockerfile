FROM jbndlr/go-example-service:0.0.1-dev AS builder
COPY jbndlr/example /root/go/src/jbndlr/example
RUN cd /root/go/src/jbndlr/example \
 && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -tags netgo -a -v -o service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /root/go/src/jbndlr/example/service .
COPY --from=builder /root/go/src/jbndlr/example/conf/defaults.yaml ./conf/defaults.yaml
CMD ["./service"]
