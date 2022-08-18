FROM golang:1.19-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o breitbandmessung

FROM zenika/alpine-chrome:latest

COPY --from=builder /build/breitbandmessung /usr/bin/breitbandmessung

ENTRYPOINT ["breitbandmessung"]
CMD ["snapshot"]
