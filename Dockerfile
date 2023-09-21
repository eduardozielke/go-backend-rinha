FROM golang:alpine as builder

RUN apk update 
RUN apk upgrade --update-cache --available
RUN apk add git make curl perl bash build-base zlib-dev ucl-dev 

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]
