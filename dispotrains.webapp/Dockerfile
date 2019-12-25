FROM node:13-buster
MAINTAINER Etienne J. Membrives <etienne@membrives.fr>

ENV TZ Europe/Paris
RUN cp /usr/share/zoneinfo/Europe/Paris /etc/localtime

COPY docker/etc/apt /etc/apt

RUN apt-get update
RUN apt-get install -y --no-install-recommends wget gnupg
RUN sh -c 'wget -qO - https://www.mongodb.org/static/pgp/server-4.2.asc | apt-key add -'

COPY docker/etc/mongodb/mongodb.list /etc/apt/sources.list.d/mongodb.list

RUN apt-get update
RUN apt-get dist-upgrade -y
RUN apt-get update
RUN apt-get install -y --no-install-recommends golang-go \
                     git tar bzip2 \
                     cron supervisor \
                     mongodb-org-shell mongodb-org-tools
RUN apt-get install -t buster-backports -y --no-install-recommends golang-go

COPY docker/etc /etc
RUN chmod 0644 /etc/crontab
COPY docker/scripts /scripts

ENV GOPATH /go
COPY src /go/src/github.com/emembrives/dispotrains/dispotrains.webapp/src
WORKDIR /go/src/github.com/emembrives/dispotrains/dispotrains.webapp/src/
RUN make
RUN ln -s /go/src/github.com/emembrives/dispotrains/dispotrains.webapp/ /dispotrains

EXPOSE 9000
WORKDIR /dispotrains/build/
CMD ["/usr/bin/supervisord"]
