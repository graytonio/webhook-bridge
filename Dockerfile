FROM alpine:3.18
ENTRYPOINT ["/webhook-bridge"]
COPY webhook-bridge /
