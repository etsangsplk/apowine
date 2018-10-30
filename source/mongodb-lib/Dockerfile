FROM alpine

RUN apk update && apk --no-cache add curl

ADD client /client

ENTRYPOINT ["/client"]
