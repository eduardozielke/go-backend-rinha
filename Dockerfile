FROM golang:alpine as builder

ARG DB_HOST=db
ARG DB_NAME=zielke_mongodb
ARG DB_USER=zielke
ARG DB_PWD=123
ARG DB_COLLECTION=pessoas
ARG SERVER_PORT=8080
ARG GIN_MODE=release

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
COPY .env /app/

EXPOSE ${SERVER_PORT}

ENV DB_HOST=${DB_HOST}
ENV DB_NAME=${DB_NAME}
ENV DB_USER=${DB_USER}
ENV DB_PWD=${DB_PWD}
ENV DB_COLLECTION=${DB_COLLECTION}
ENV GIN_MODE=${GIN_MODE}

CMD ["./main"]
