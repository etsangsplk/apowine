FROM alpine

RUN apk update && apk --no-cache add curl

RUN mkdir -p /apowine

ADD client /apowine/client

ENTRYPOINT ["/apowine/client"]
