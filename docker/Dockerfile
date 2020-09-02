FROM registry.cn-hangzhou.aliyuncs.com/hongweigao/rashomon_base:0.0.1

ENV GOPATH=/root/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
ENV GO111MODULE=on
ENV GOPROXY=https://mirrors.aliyun.com/goproxy/

RUN mkdir -p ${GOPATH}/src/github.com/hongweikkx/rashomon
COPY ./ ${GOPATH}/src/github.com/hongweikkx/rashomon
WORKDIR ${GOPATH}/src/github.com/hongweikkx/rashomon

RUN make build