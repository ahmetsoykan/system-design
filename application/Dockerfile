# BUILD
FROM golang:1.20 as builder

WORKDIR /go/url-shortener
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o ./dist/urls  cmd/main.go

## -------------------------------------------------------------------------------- ##
# DEPLOY
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=builder /go/url-shortener/dist/urls /urls

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/urls"]