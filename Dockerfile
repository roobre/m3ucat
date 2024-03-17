FROM golang:1.22-alpine as builder

WORKDIR /m3ucat

COPY . .
RUN go build -o /bin/m3ucat ./cmd

FROM alpine:3.19.1

RUN apk update && apk add bash
COPY --from=builder /bin/m3ucat /usr/local/bin/
ENTRYPOINT [ "/usr/local/bin/m3ucat" ]
CMD ["-"]
