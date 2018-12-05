FROM golang
LABEL maintainer="heqingliang"
RUN go get -v github.com/garyburd/redigo/redis
COPY engine /go/bin/
COPY wait-for-it.sh /
CMD ["engine"]