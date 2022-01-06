ARG GOPROXY=https://goproxy.cn,direct

# STEP 1 build executable binary from alpine:golang

FROM golang:1.16-alpine AS builder
LABEL maintainer="EthanWang <git@ethan.wang>"

ENV GOPROXY=${GOPROXY}

RUN apk add --update --no-cache build-base musl-dev linux-headers git \
    pkgconf g++ gmp-dev openssl-dev

WORKDIR /src

ADD ./ /src

RUN go mod download
WORKDIR /src/price_server

RUN make warehouse

FROM alpine:latest
LABEL maintainer="EthanWang <git@ethan.wang>"

RUN apk add --update --no-cache tzdata ca-certificates bash py-pip groff \
    musl gmp openssl libstdc++ fish shadow \
    && usermod -s `which fish` root

ENV DEBIAN_FRONTEND=noninteractive

COPY --from=builder /src/price_server/build/warehouse /opt/warehouse
COPY --from=builder /src/price_server/configs /opt/configs


EXPOSE 5566

WORKDIR /opt
#CMD ["/opt/warehouse"]
