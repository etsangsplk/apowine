FROM alpine
MAINTAINER Aporeto Inc <info@aporeto.com>

RUN apk update && apk --no-cache add curl

ADD client /client

ENTRYPOINT ["/client"]
