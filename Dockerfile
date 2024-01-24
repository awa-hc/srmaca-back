FROM golang:1.21.6 as builder

WORKDIR /app

COPY . .

RUN go mod download \
    && go mod tidy \
    && go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

USER nobody
EXPOSE 8080

CMD [ "./main" ]