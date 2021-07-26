FROM golang:1.16-stretch as builder
ENV  GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /generator

FROM debian:stretch-slim

COPY --from=builder /generator /generator

ENTRYPOINT [ "./generator" ]
CMD [ "--collectorUrl", "ingest.lightstep.com:443", "--paramsFile", "/topo.json" ]