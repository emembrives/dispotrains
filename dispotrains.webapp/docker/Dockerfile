FROM golang:latest
MAINTAINER Etienne J. Membrives <etienne@membrives.fr>

ENV TZ Europe/Paris
RUN cp /usr/share/zoneinfo/Europe/Paris /etc/localtime

RUN apt-get update
RUN apt-get dist-upgrade -y

COPY src /dispotrains/src/
WORKDIR /dispotrains/src/
RUN make
