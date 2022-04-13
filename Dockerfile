FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /build

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o api_transformer main.go

FROM alpine:latest as prod

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=0 /build/api_transformer .

EXPOSE 8080
CMD ["./api_transformer"]