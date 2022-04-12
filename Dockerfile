FROM registry.gitlab.com/etke.cc/base AS builder

WORKDIR /buscarron
COPY . .
RUN make build

FROM alpine:latest

ENV BUSCARRON_DB_DSN /data/buscarron.db

RUN apk --no-cache add ca-certificates tzdata olm && update-ca-certificates && \
    adduser -D -g '' buscarron && \
    mkdir /data && chown -R buscarron /data

COPY --from=builder /buscarron/buscarron /opt/buscarron/buscarron

WORKDIR /opt/buscarron
USER buscarron

ENTRYPOINT ["/opt/buscarron/buscarron"]

