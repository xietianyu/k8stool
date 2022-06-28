FROM golang:alpine AS builder

ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /home/tyler

COPY . .

RUN go mod download
RUN go build -o workflowapi ./main.go
EXPOSE 8888
CMD ["./workflowapi"]