FROM golang:1.20-alpine  as builder

#RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go get -u ./...

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./Entities ./Entities
COPY ./pkg ./pkg

#RUN mkdir "logs"
#RUN mkdir "metrics"
#RUN mkdir "postgres_data"
#COPY ./queries .
#RUN mkdir "redis_data"
#RUN mkdir "static"

COPY ./static ./static

RUN touch ./cmd/app/main.go
RUN GOOS=linux GOARCH=amd64 go build -o ./.bin/app ./cmd/app/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/.bin/app .
#RUN mkdir "static"

EXPOSE 4040
EXPOSE 9091

RUN apk add dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["./app"]