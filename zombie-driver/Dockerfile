FROM golang

RUN go get -u github.com/golang/dep/...

ARG HTTP_ADDR

ENV HTTP_ADDR ${HTTP_ADDR}
ENV SRC_DIR=/go/src/github.com/spirosoik/go-driver-microservices/zombie-driver

WORKDIR $SRC_DIR
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

ADD . $SRC_DIR
RUN mkdir /app
RUN go build -o ./bin/zombie-driver -v cmd/zombie-driver/main.go; cp ./bin/zombie-driver /app/zombie-driver

WORKDIR /app/

ENTRYPOINT ["sh", "-c", "/app/zombie-driver --http.addr=${HTTP_ADDR}"]