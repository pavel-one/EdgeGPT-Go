FROM golang:1.20 as builder

LABEL org.opencontainers.image.source=https://github.com/pavel-one/EdgeGPT-Go

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build -o ./build cmd/grpc.go

RUN rm -rf internal/ cmd/ config/ grpc/ pkg/

CMD ["/app/build"]
