FROM node:20-bookworm AS builder

RUN apt update
RUN apt full-upgrade -y
RUN apt update
RUN apt install -y --no-install-recommends golang-go librocksdb-dev

ENV GOPATH /go
COPY src /go/src/github.com/emembrives/dispotrains/dispotrains.webapp/src
WORKDIR /go/src/github.com/emembrives/dispotrains/dispotrains.webapp/src/
RUN make

FROM bitnami/minideb:latest AS runner
ENV TZ Europe/Paris
RUN cp /usr/share/zoneinfo/Europe/Paris /etc/localtime
RUN install_packages librocksdb7.8 supervisor curl cron ca-certificates

COPY --from=builder /go/src/github.com/emembrives/dispotrains/dispotrains.webapp/build /dispotrains

COPY docker/etc /etc
RUN chmod 0644 /etc/crontab
COPY docker/scripts /scripts

EXPOSE 9000
WORKDIR /dispotrains/
CMD ["/usr/bin/supervisord"]
