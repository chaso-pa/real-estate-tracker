FROM golang:1.25 AS build
WORKDIR /go/src
COPY cmd ./cmd
COPY internal ./internal
COPY go.sum .
COPY go.mod .

ENV CGO_ENABLED=0

RUN go build -C /go/src/cmd/server/ -o /go/src/server .

FROM alpine:3.20 AS runtime
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
ENV GIN_MODE=release
COPY --from=build /go/src/server ./

EXPOSE 80/tcp

ENV PORT 80
ENTRYPOINT ["./server"]
