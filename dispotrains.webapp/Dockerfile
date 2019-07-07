FROM node:11-stretch
MAINTAINER Etienne J. Membrives <etienne@membrives.fr>

ENV TZ Europe/Paris
RUN cp /usr/share/zoneinfo/Europe/Paris /etc/localtime

COPY docker/etc/apt /etc/apt

RUN apt-get update
RUN apt-get dist-upgrade -y
RUN apt-get update
RUN apt-get install -y --no-install-recommends golang-go \
                     git tar bzip2 \
                     mongodb-clients mongo-tools \
                     cron supervisor
RUN apt-get install -t stretch-backports -y --no-install-recommends golang-go

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
