FROM golang:alpine as builder

LABEL maintainer "Jon Friesen <jon@jonfriesen.ca>"

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

WORKDIR /go/src/github.com/jonfriesen/pglint
COPY . .

RUN go build -o bin/pglint-cli cmd/pglint-cli/main.go

FROM alpine:latest

RUN apk --no-cache add postgresql-dev

COPY --from=builder /go/src/github.com/jonfriesen/pglint/bin/pglint-cli /usr/bin/pglint-cli

# any commands added after will be appended
ENTRYPOINT [ "pglint-cli" ]
