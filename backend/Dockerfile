FROM golang:alpine as builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /home
COPY . .

RUN go mod download
RUN go build -o course_design main.go

FROM alpine as runner
ENV WORKDIR=/home
WORKDIR $WORKDIR

RUN apk add tzdata --no-cache \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

COPY --from=builder $WORKDIR/tce_label  $WORKDIR
COPY --from=builder $WORKDIR/config  $WORKDIR/config
RUN mkdir -p /home/fileTmp

CMD ["./course_design"]
