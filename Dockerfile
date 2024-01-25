FROM golang:alpine AS builder
WORKDIR /app
COPY . . 
RUN go build -o backend

FROM alpine 
WORKDIR /app
COPY --from=builder /app/backend .
CMD ["./backend"]