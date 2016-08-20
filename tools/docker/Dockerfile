FROM golang:latest

RUN go get -u v2ray.com/core/... \
  && rm -f $GOPATH/bin/build \
  && go install v2ray.com/core/tools/build \
  && build \
  && mv bin/v2ray-custom-linux-64/v2ray bin \
  && rm -rf pkg src

EXPOSE 27183
COPY gen-server-cfg.sh /go/gen-server-cfg.sh
COPY docker-entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["v2ray", "--config=/go/server-cfg.json"]
