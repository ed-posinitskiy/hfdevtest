FROM golang:1.17.10-alpine3.16 as build

ENV WORKDIR=/go/src/recipe-collector

WORKDIR "${WORKDIR}"

RUN apk update && apk add --no-cache \
        git

COPY . .

RUN CGO_ENABLED=0 go test ./... -v \
 && GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/collector

FROM alpine:3.11.3 as runtime

COPY --from=build /go/bin/collector /usr/local/bin/collector

ENTRYPOINT ["collector"]