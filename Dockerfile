FROM golang:1.16.6-alpine AS build-env

WORKDIR /build

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk --no-cache add git=~2 upx=~3

COPY *.go go.mod go.sum /build/

RUN go version
RUN go build -tags docker -ldflags '-s -w'

FROM scratch

EXPOSE 80

ENV PATH=/

COPY --from=build-env /build/aula-assistant /aula-assistant

ENTRYPOINT ["aula-assistant"]
