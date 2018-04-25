FROM alpine

MAINTAINER Aporeto Inc <info@aporeto.com>

RUN apk update && apk --no-cache add curl

RUN mkdir -p /apowine

ADD producerbeer /apowine/producerbeer

ENTRYPOINT ["/apowine/producerbeer"]
