FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .
COPY . .
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && go mod tidy && go build -o user-service ./cmd/main.go

FROM alpine

WORKDIR /build
COPY --from=builder /build /build

CMD ["./user-service"]