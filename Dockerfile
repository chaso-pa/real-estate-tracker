FROM golang:1.23 AS build
WORKDIR /go/src
COPY cmd ./cmd
COPY internal ./internal
COPY go.sum .
COPY go.mod .

ENV CGO_ENABLED=0

RUN go build -C cmd/server/  -o server .

FROM alpine:3.20 AS runtime
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
ENV GIN_MODE=release
COPY --from=build /go/src/server ./

EXPOSE 80/tcp

ENV PORT 80
ENTRYPOINT ["./server"]
