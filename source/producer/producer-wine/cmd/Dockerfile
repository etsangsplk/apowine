FROM alpine

RUN apk update && apk --no-cache add curl

RUN mkdir -p /apowine

ADD producerwine /apowine/producerwine

ADD producerwine.txt /apowine/producerwine.txt

ENTRYPOINT ["/apowine/producerwine"]
