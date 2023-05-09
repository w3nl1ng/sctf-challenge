FROM golang:latest AS builder

RUN mkdir /app
ADD ./ /app
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o main ./app/server


FROM alpine:latest AS production
COPY --from=builder /app/main .
CMD ["./main"]⏎ 