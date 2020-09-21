FROM golang:1.14.7-alpine

ENV APP_DIR $GOPATH/src/github.com/hoshiin/csv-inserter
ENV PATH $GOPATH/bin:$PATH
ENV GO111MODULE on

RUN apk update && \
	apk add git make gcc g++

# TODO: realize cannot install when GO111MODULE on
# https://github.com/oxequa/realize/issues/253#issuecomment-534826718
RUN GO111MODULE=off && go get -u github.com/oxequa/realize
RUN GO111MODULE=off go get -u github.com/go-delve/delve/cmd/dlv
RUN go get -u golang.org/x/tools/cmd/goimports

ADD . $APP_DIR
WORKDIR $APP_DIR

RUN apk --no-cache add tzdata curl && \
	cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/csv-inserter  -ldflags="-s -w"
