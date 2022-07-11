#Stage I Building Binary
FROM golang:1.18-alpine AS builder
RUN apk add build-base
RUN mkdir app
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/main .

#Stage II
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .
EXPOSE 8050
CMD [ "./main"]
