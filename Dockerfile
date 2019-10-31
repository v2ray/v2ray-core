############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git bash wget
WORKDIR /go/src/v2ray.com/core
RUN git clone --progress --depth=1 --branch dockerdev https://github.com/v2fly/v2ray-core.git . && \
    bash -x ./release/user-package.sh nosource noconf abpathtgz=/tmp/v2ray.tgz
############################
# STEP 2 build a small image
############################
FROM alpine
COPY --from=builder /tmp/v2ray.tgz /tmp
RUN apk update && apk add ca-certificates && \
    mkdir -p /usr/bin/v2ray && \
    tar xvfz /tmp/v2ray.tgz -C /usr/bin/v2ray
ENTRYPOINT ["/usr/bin/v2ray/v2ray"]
