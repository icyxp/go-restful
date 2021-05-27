# build base for boost speed
FROM golang:1.13.8 as build_go-restful_base

WORKDIR /go/src/go-restful

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

COPY go.mod .
COPY go.sum .

RUN go get github.com/google/gops && \
    go mod download

# Generate executable file
FROM build_go-restful_base AS go_builder

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -a -v -ldflags="-s" ./ 

# final image
FROM icyboy/centos:7.6-base

WORKDIR /opt

COPY --from=go_builder /go/bin/go-restful /opt/go-restful
COPY --from=go_builder /go/bin/gops /bin/gops

ENTRYPOINT [ "/opt/go-restful" ]
