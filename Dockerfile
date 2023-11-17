ARG CMD=app-main

FROM golang:latest AS go-builder
ARG CMD
ENV CGO_ENABLED 0
WORKDIR /src
ADD . .
RUN go build --buildvcs=false -o $CMD

FROM alpine:latest
ARG CMD
RUN apk update
COPY --from=go-builder /src/$CMD /$CMD

RUN ln -s /$CMD /entrypoint

CMD ["/entrypoint"]