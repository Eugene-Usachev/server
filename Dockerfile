FROM golang:1.20-alpine as builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go get -u ./...

COPY . .

RUN touch ./cmd/app/main.go
RUN GOOS=linux GOARCH=amd64 go build -o ./.bin/app ./cmd/app/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/.bin/app .
COPY static/ /root/static/

EXPOSE 4040

RUN apk add dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["./app"]