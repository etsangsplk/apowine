FROM alpine

RUN apk update && apk --no-cache add curl

RUN mkdir -p /apowine

ADD server /apowine/server

EXPOSE 3000

ENTRYPOINT ["/apowine/server"]
