# build base for boost speed
FROM golang:1.13.8 as build_go-restful_base

WORKDIR /go/src/go-restful

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache upx ca-certificates tzdata

COPY go.mod .
COPY go.sum .

RUN go get github.com/google/gops && \
    go mod download

# Generate executable file
FROM build_go-restful_base AS go_builder

COPY . .

# -s 去掉符号表信息, panic时候的stack trace就没有任何文件名/行号信息了.
# -w 去掉DWARF调试信息，得到的程序就不能用gdb调试了
# s/w 不要同时使用,w中包含s
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o go-restful && \
    upx --best --lzma go-restful -o _upx_go-restful && \
    mv -f _upx_go-restful go-restful

# final image
FROM icyboy/centos:7.6-base

WORKDIR /opt

COPY --from=go_builder /go/bin/go-restful /opt/go-restful
COPY --from=go_builder /go/bin/gops /bin/gops

ENTRYPOINT [ "/opt/go-restful" ]
