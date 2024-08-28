FROM golang:alpine AS builder
RUN apk add -U --no-cache ca-certificates && update-ca-certificates

ENV GO111MODULE=on\
    CGO_ENABLED=0\
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum /cmd/server/main.go ./
RUN go mod download

COPY . ./
RUN go build -o main .

WORKDIR /dist
RUN cp /build/main .
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /dist/main .

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip
ENV TZ=Asia/Seoul

ENTRYPOINT [ "./main", "run" ]