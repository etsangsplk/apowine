FROM alpine

MAINTAINER Aporeto Inc <info@aporeto.com>

RUN apk update && apk --no-cache add curl && apk add --no-cache bash

RUN mkdir -p /apowine

ADD frontend /apowine/frontend

EXPOSE 3005

ENTRYPOINT ["/apowine/frontend"]
