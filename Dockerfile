FROM golang:1.8

LABEL maintainer 'nsosnov@dataart.com'

WORKDIR /go/src/bidder
COPY . .

RUN go-wrapper download \
    && go-wrapper install

EXPOSE 8080

CMD ["/go/bin/bidder"]
