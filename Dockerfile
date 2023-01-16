FROM registry.gitlab.com/etke.cc/base/build AS builder

WORKDIR /buscarron
COPY . .
RUN just build

FROM registry.gitlab.com/etke.cc/base/app

ENV BUSCARRON_DB_DSN /data/buscarron.db

COPY --from=builder /buscarron/buscarron /bin/buscarron

USER app

ENTRYPOINT ["/bin/buscarron"]

