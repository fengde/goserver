FROM golang:1.18 AS builder

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod vendor -o app main.go

FROM centos:7.9.2009

WORKDIR /runtime
COPY --from=builder /build/app /runtime/app
COPY --from=builder /build/conf/config.yaml /runtime/conf/config.yaml

EXPOSE 8888

ENV TZ Asia/Shanghai
RUN ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone

ENTRYPOINT [ "sh", "-c", "/runtime/app" ]