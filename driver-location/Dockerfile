FROM golang

RUN go get -u github.com/golang/dep/...

ARG NSQ_LOOKUP_ADDR
ARG REDIS_ADDR
ARG HTTP_ADDR

ENV NSQ_LOOKUP_ADDR ${NSQ_LOOKUP_ADDR}
ENV HTTP_ADDR ${HTTP_ADDR}
ENV REDIS_ADDR ${REDIS_ADDR}
ENV SRC_DIR=/go/src/github.com/spirosoik/go-driver-microservices/driver-location

WORKDIR $SRC_DIR
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

ADD . $SRC_DIR
RUN mkdir /app
RUN go build -o ./bin/driver-location -v cmd/driver-location/main.go; cp ./bin/driver-location /app/driver-location

WORKDIR /app/

ENTRYPOINT ["sh", "-c", "/app/driver-location --redis.addr=${REDIS_ADDR} --http.addr=${HTTP_ADDR} --lookup.addr=${NSQ_LOOKUP_ADDR}"]