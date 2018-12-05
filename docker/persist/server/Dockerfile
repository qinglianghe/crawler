FROM golang
LABEL maintainer="heqingliang"
RUN go get -v github.com/olivere/elastic
COPY itemsaver /go/bin/
COPY wait-for-it.sh /
CMD ["itemsaver"]