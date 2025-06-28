FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o payment-service main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=build /app/payment-service .
COPY --from=build /app .
EXPOSE 5001
CMD ["./payment-service"] 