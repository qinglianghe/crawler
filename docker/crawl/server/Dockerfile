FROM golang
LABEL maintainer="heqingliang"
RUN go get -v -u github.com/gpmgo/gopm && \ 
    gopm get -g -v golang.org/x/text && \
    gopm get -g -v golang.org/x/net/html
COPY crawl /go/bin/
CMD ["crawl"]
