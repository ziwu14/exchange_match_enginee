FROM golang:1.8

# RUN echo 'sysctl -w net.core.somaxconn=65535' > /etc/rc.local && echo 'vm.overcommit_memory = 1' > /etc/sysctl.conf && echo 'echo never > /sys/kernel/mm/transparent_hugepage/enabled' > /etc/rc.local
ENV GOPATH=/go
ENV COMPOSE_HTTP_TIMEOUT=3600

COPY ./src /go/src/app/

WORKDIR /go/src/app/

# All packages from github are in local project
# IF not, write go get command here.
RUN go get "github.com/gomodule/redigo/redis"
RUN go get "github.com/sirupsen/logrus"
RUN go get "github.com/beevik/etree"