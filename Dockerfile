FROM golang:alpine AS builder

ADD . /go/src/videocvt/
RUN apk add git && \
    go get -u github.com/gin-gonic/gin && \
    go get -u github.com/sirupsen/logrus && \
    cd /go/src/videocvt/ && \
    go build -o videocvt

FROM jrottenberg/ffmpeg:3.2-alpine
COPY --from=builder /go/src/videocvt/videocvt /usr/bin/videocvt
CMD ["11111111"]
ENTRYPOINT ["sleep"]



