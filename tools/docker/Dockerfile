FROM golang:latest

RUN go get -u github.com/v2ray/v2ray-core
RUN rm -f $GOPATH/bin/build
RUN go install github.com/v2ray/v2ray-core/tools/build
RUN $GOPATH/bin/build

EXPOSE 27183
ADD server-cfg.json /go/server-cfg.json
CMD /go/bin/v2ray-custom-linux-64/v2ray --config="/go/server-cfg.json"
