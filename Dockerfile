FROM golang:1.15 AS builder

ARG GOBIN=/go/bin/
ARG GOOS=linux
ARG GOARCH=amd64
ARG GOPATH=$HOME/go
ARG CGO_ENABLED=0
ARG PKG_NAME=github.com/AccumulateNetwork/accumulated
ARG PKG_PATH=${GOPATH}/src/${PKG_NAME}

WORKDIR ${PKG_PATH}
COPY . ${PKG_PATH}/

#RUN go mod download
RUN go run ./cmd/accumulated init -n "Badlands" --relay-to "Badlands"
RUN go build -o /go/bin/accumulated ./cmd/accumulated/

FROM alpine:3.7

RUN set -xe && \
  apk --no-cache add bash ca-certificates inotify-tools && \
  addgroup -g 1000 app && \
  adduser -D -G app -u 1000 app

WORKDIR /home/app

COPY --from=builder /go/bin/accumulated ./
COPY --from=builder /root/.accumulate ./.accumulate

##
#COPY ./docker-entrypoint.sh ./docker-entrypoint.sh
#RUN chmod 775 docker-entrypoint.sh
#ENTRYPOINT ["docker-entrypoint.sh"]
#CMD run

RUN \
  mkdir ./values && \
  chown -R app:app /home/app

USER app

EXPOSE 35554 35554

CMD [ "./accumulated", "run", "-n", "0" ]
