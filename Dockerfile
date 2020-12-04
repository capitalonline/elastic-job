FROM golang:1.12.9 as builder
ENV GOPROXY https://goproxy.io
#ENV GO111MODULE on

WORKDIR /go/cache

ADD go.mod .
RUN go mod download
WORKDIR /go/release
ADD . .
RUN  go build -o mongodb-job

FROM centos:centos7.6.1810 as pro
#RUN yum update -y
COPY get-pip.py .
COPY ansible.cfg /etc/ansible/ansible.cfg
RUN python get-pip.py;
RUN yum install -y sshpass openssl openssl-devel openssh-clients
RUN mkdir -p /script
RUN pip install  ansible==2.9.0 -i https://pypi.douban.com/simple
RUN mkdir -p /app/logs
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /go/release/mongodb-job /app
COPY --from=builder /go/release/script /app/script
WORKDIR /app
CMD ["./mongodb-job",""]
