FROM golang:latest

EXPOSE 8080

RUN mkdir /app
COPY ./main.go /app

WORKDIR /app
ENV GOPROXY "https://goproxy.cn,direct"
RUN go mod init server
RUN go mod tidy
RUN go build -o main .

CMD [ "./main" ]