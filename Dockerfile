FROM golang:1.16 AS build_base

WORKDIR /go/src

# Force the go compiler to use modules
ENV GO111MODULE=on

COPY go.mod .

#COPY go.sum .

RUN go mod download

# This image builds the weavaite server
FROM build_base AS builder

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o instagram .

FROM alpine

COPY --from=builder /go/src/instagram /app/

COPY --from=builder /go/src/instagram.slots /app/

RUN apk add --no-cache curl

WORKDIR /app

CMD ["./instagram"]
