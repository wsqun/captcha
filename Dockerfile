FROM golang:1.17 as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN  go mod download -x

COPY . .
RUN make build
#RUN make buildCompare

FROM alpine:3.9.5 as captcha
WORKDIR /app
COPY --from=builder /usr/src/app/captcha .
COPY --from=builder /usr/src/app/examples/comic.ttf .
RUN  mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
EXPOSE 8085
ENTRYPOINT ["./captcha"]