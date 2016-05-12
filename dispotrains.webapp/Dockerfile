FROM golang:latest
MAINTAINER Etienne J. Membrives <etienne@membrives.fr>

ENV TZ Europe/Paris
RUN cp /usr/share/zoneinfo/Europe/Paris /etc/localtime

RUN apt-get update
RUN apt-get dist-upgrade -y
RUN apt-get install -y nodejs-legacy npm git
RUN apt-get install -y mongodb-clients tar bzip2
RUN apt-get install -y cron supervisor
RUN npm install -g bower

COPY docker/etc /etc
RUN chmod 644 /etc/crontab

COPY src /go/src/github.com/emembrives/dispotrains/dispotrains.webapp/src
WORKDIR /go/src/github.com/emembrives/dispotrains/dispotrains.webapp/src/
RUN make
RUN ln -s /go/src/github.com/emembrives/dispotrains/dispotrains.webapp/ /dispotrains

EXPOSE 9000
CMD ["/usr/bin/supervisord"]
