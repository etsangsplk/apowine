FROM alpine

MAINTAINER Aporeto Inc <info@aporeto.com>

RUN apk update && apk --no-cache add curl

RUN mkdir -p /apowine

ADD server /apowine/server

EXPOSE 3000 5001

ENTRYPOINT ["/apowine/server"]
