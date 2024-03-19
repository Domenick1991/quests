FROM golang:1.22.1-alpine3.19 as builder
WORKDIR /usr/local/src
COPY ["go.mod", "go.sum", "./"]
RUN go mod download
COPY . ./
RUN go build -o ./bin/app main.go

FROM alpine:3.19 as app
COPY --from=builder /usr/local/src/bin/app /
COPY config.yml config.yml
CMD ["/app"]