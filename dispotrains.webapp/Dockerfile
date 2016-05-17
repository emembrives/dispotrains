FROM golang:latest
MAINTAINER Etienne J. Membrives <etienne@membrives.fr>

ENV TZ Europe/Paris
RUN cp /usr/share/zoneinfo/Europe/Paris /etc/localtime

RUN apt-get update
RUN apt-get dist-upgrade -y
RUN apt-get install -y nodejs-legacy npm \
                     git tar bzip2 \
                     mongodb-clients \
                     cron supervisor
RUN npm install -g bower

COPY docker/etc /etc
RUN chmod 0644 /etc/crontab
COPY docker/scripts /scripts

COPY src /go/src/github.com/emembrives/dispotrains/dispotrains.webapp/src
WORKDIR /go/src/github.com/emembrives/dispotrains/dispotrains.webapp/src/
RUN make
RUN ln -s /go/src/github.com/emembrives/dispotrains/dispotrains.webapp/ /dispotrains

EXPOSE 9000
WORKDIR /dispotrains/build/
CMD ["/usr/bin/supervisord"]
