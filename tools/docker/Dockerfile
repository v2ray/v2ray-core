FROM phusion/baseimage

ENV GOPATH /v2ray
ENV PATH $PATH:/usr/local/go/bin

RUN apt-get update
RUN apt-get install -y git
RUN apt-get clean

RUN mkdir -p /v2ray/logs
RUN curl -o /go.tar.gz https://storage.googleapis.com/golang/go1.5.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf /go.tar.gz
RUN rm /go.tar.gz


RUN go get -u github.com/v2ray/v2ray-core
RUN rm -f $GOPATH/bin/build
RUN go install github.com/v2ray/v2ray-core/tools/build
RUN $GOPATH/bin/build

EXPOSE 27183
ADD server-cfg.json /v2ray/server-cfg.json
CMD /v2ray/bin/v2ray-custom-linux-64/v2ray --config="/v2ray/server-cfg.json"
