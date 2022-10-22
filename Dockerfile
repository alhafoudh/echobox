FROM --platform=linux/amd64 golang:1.19.2 AS builder
WORKDIR /go/src/github.com/alhafoudh/echobox/
COPY main.go go.mod go.sum ./
RUN go build -a -o echobox .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/alhafoudh/echobox/echobox /app/echobox
CMD ["/app/echobox"]