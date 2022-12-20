FROM golang:1.16 AS builder
WORKDIR /go/src/github.com/yokawasa/awsping/
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY . .
RUN GO111MODULE=on go build -o awsping

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/yokawasa/awsping/awsping ./
CMD ["./awsping"]
