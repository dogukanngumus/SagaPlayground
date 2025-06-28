FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o stock-service main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=build /app/stock-service .
COPY --from=build /app .
EXPOSE 5002
CMD ["./stock-service"] 