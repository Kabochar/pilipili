FROM golang as build

# 设置 Proxy
ENV GOPROXY=https://goproxy.io

ADD . /pilipili

WORKDIR /pilipili

# 交叉编译
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api_server

FROM alpine:3.7

ENV REDIS_ADDR=""
ENV REDIS_PW=""
ENV REDIS_DB=""
ENV MysqlDSN=""
ENV GIN_MODE="release"
ENV PORT=3000

# 设置软件源
RUN echo "http://mirrors.aliyun.com/alpine/v3.7/main/" > /etc/apk/repositories && \
    apk update && \
    apk add ca-certificates && \
    echo "hosts: files dns" > /etc/nsswitch.conf && \
    mkdir -p /www/conf

WORKDIR /www

COPY --from=build /pilipili/api_server /usr/bin/api_server
ADD ./conf /www/conf

# 加权限
RUN chmod +x /usr/bin/api_server

ENTRYPOINT ["api_server"]