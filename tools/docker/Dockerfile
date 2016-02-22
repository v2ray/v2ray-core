FROM golang:latest

RUN go get -u github.com/v2ray/v2ray-core \
  && rm -f $GOPATH/bin/build \
  && go install github.com/v2ray/v2ray-core/tools/build \
  && build \
  && rm -rf pkg src

EXPOSE 27183
COPY gen-server-cfg.sh /go/gen-server-cfg.sh
COPY docker-entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/go/bin/v2ray-custom-linux-64/v2ray", "--config=/go/server-cfg.json"]
