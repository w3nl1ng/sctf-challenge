FROM golang:latest

EXPOSE 8080

RUN mkdir /app
COPY ./ /app

WORKDIR /app
ENV GOPROXY "https://goproxy.cn,direct"
RUN go build -o main .

CMD [ "./main", "-h", "172.17.0.1", "-P", "3307", "-n", "sctf_db", "-u", "sctf", "-p", "sctf2023recv!!!"]