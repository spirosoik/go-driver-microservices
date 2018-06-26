FROM golang

RUN go get -u github.com/golang/dep/...

ARG NSQ_ADDR
ARG HTTP_ADDR

ENV NSQ_ADDR ${NSQ_ADDR}
ENV HTTP_ADDR ${HTTP_ADDR}
ENV SRC_DIR=/go/src/github.com/spirosoik/go-driver-microservices/gateway

WORKDIR $SRC_DIR
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

ADD . $SRC_DIR
RUN mkdir /app
RUN go build -o ./bin/gateway -v cmd/gateway/main.go; cp ./bin/gateway /app/gateway

WORKDIR /app/

ENTRYPOINT ["sh", "-c", "/app/gateway --config=$SRC_DIR/config.yaml --http.addr=${HTTP_ADDR} --nsq.addr=${NSQ_ADDR}"]