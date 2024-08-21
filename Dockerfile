FROM ghcr.io/etkecc/base/build AS builder

WORKDIR /app
COPY . .
RUN just build

FROM scratch

ENV BUSCARRON_DB_DSN /data/buscarron.db

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/PROJECT /bin/buscarron

USER app

ENTRYPOINT ["/bin/buscarron"]

