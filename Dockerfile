FROM alpine:latest AS cert-provider
RUN apk --no-cache add ca-certificates

FROM scratch

COPY --from=cert-provider /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY bin/toygoproxy /toygoproxy

ENTRYPOINT ["/toygoproxy"]
