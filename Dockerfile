FROM golang:1.20-alpine as builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go get -u ./...

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o ./.bin/app ./cmd/app/main.go

# Final stage
FROM alpine:latest

ENV DB_HOST=localhost \
    DB_PASSWORD=db_create.go \
    DB_USERNAME=postgres \
    DB_PORT=5436 \
    PORT=4040 \
    DB_NAME=hey \
    SSL_MODE=disable \
    CORES=5 \
    SALT=uipads0797Wy \
    JWT_SECRET_KEY=pz1xvp44cMoca1dmf9HJfa \
    JWT_SECRET_KEY_FOR_LONGLIVE_TOKEN=zp1mvn65cmOca1hgf9HJfa \
    MAX_FILE_SIZE=66060288 \
    GIN_MODE=release
# todo remove  GIN_MODE=release

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/.bin/app .
COPY static/ /root/static/

EXPOSE 4040

CMD ["./app"]