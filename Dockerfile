FROM golang:1.17.3-alpine3.14 as build-env
LABEL maintainer="ForestVPN.com <support@forestvpn.com>"

RUN apk add --no-cache git gcc
RUN mkdir /app

WORKDIR /app

COPY . .

RUN export appVersion=$(git describe --tags `git rev-list -1 HEAD`) && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags "-X main.appVersion=$appVersion" \
      -o wgrest cmd/wgrest-server/main.go

FROM alpine:3.14
LABEL maintainer="ForestVPN.com <support@forestvpn.com>"

COPY --from=build-env /app/wgrest .

EXPOSE 8080/tcp

USER 1001

ENTRYPOINT ["./wgrest"]