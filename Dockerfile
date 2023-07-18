FROM golang:1.16.5-alpine AS builder

RUN apk add netcat-openbsd ca-certificates
WORKDIR /go/src/github.com/and3rson/dunai
COPY go.mod go.sum ./
RUN go mod download -x
COPY cmd ./cmd
COPY dunai ./dunai
RUN mkdir /out && CGO_ENABLED=0 go build -o /out/dunai ./cmd/dunai/main.go
# RUN go get && mkdir /out && CGO_ENABLED=0 go build -o /out/dunai && cp -r data static templates /out
# RUN wget https://busybox.net/downloads/binaries/1.30.0-i686/busybox -O /out/busybox && chmod +x /out/busybox

FROM scratch
WORKDIR /etc/ssl/certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /root
COPY --from=builder /out/ /root/
ENTRYPOINT ["/root/dunai"]
