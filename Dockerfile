FROM golang:1.12.4-alpine3.9 as builder

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -installsuffix 'static' -o /app .

FROM alpine:3.9

COPY --from=builder /app /usr/local/bin/app

CMD ["app"]
