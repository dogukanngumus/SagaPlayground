# Saga Choreography Demo - Go & .NET Microservices

Bu proje, Saga Choreography pattern'ini kullanarak Go ve .NET mikroservislerinin asenkron event'lerle haberleştiği bir demo uygulamasıdır.

## 🏗️ Mimari

### Servisler
- **Order Service** (.NET): Sipariş oluşturma ve yönetimi
- **Payment Service** (Go): Ödeme işlemleri
- **Stock Service** (Go): Stok yönetimi

### Event Akışı
```
Order Created → Payment Requested → Payment Processed → Stock Reserved → Order Confirmed
```

### Teknolojiler
- **Message Broker**: RabbitMQ
- **Databases**: SQLite (Order), PostgreSQL (Payment, Stock)
- **Patterns**: Saga Choreography, Inbox/Outbox Pattern
- **Event Format**: JSON

## 🚀 Çalıştırma

### Gereksinimler
- Docker & Docker Compose
- .NET 8 SDK
- Go 1.21+

### Adımlar
1. Projeyi klonlayın
2. Root dizinde çalıştırın:
```bash
docker-compose up --build -d
```

3. Servisleri başlatın:
```bash
# Order Service (.NET)
cd order-service
dotnet run

# Payment Service (Go)
cd payment-service
go run main.go

# Stock Service (Go)
cd stock-service
go run main.go
```

## 📋 Örnek Akış

### 1. Sipariş Oluştur
```bash
curl -X POST http://localhost:5100/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customerId": "customer-123",
    "items": [
      {
        "productId": "product-1",
        "quantity": 2,
        "price": 25.50
      }
    ]
  }'
```

### 2. Event Akışını İzle
Logları takip ederek event akışını gözlemleyin:
- Order Service: Sipariş oluşturuldu
- Payment Service: Ödeme talebi alındı
- Stock Service: Stok rezervasyonu yapıldı
- Order Service: Sipariş onaylandı

## 🔧 Konfigürasyon

### Environment Variables
- `RABBITMQ_URL`: RabbitMQ bağlantı URL'i
- `DB_CONNECTION`: Veritabanı bağlantı string'i
- `SERVICE_PORT`: Servis port'u

### Portlar
- Order Service: 5100
- Payment Service: 5001
- Stock Service: 5002
- RabbitMQ: 5672 (AMQP), 15672 (Management UI)

## 📁 Proje Yapısı

```
SagaPlayground/
├── docker-compose.yml          # Docker ortamı
├── README.md                   # Bu dosya
├── order-service/              # .NET Order servisi
│   ├── Program.cs
│   ├── Controllers/
│   ├── Models/
│   ├── Services/
│   └── order-service.csproj
├── payment-service/            # Go Payment servisi
│   ├── main.go
│   ├── handlers/
│   ├── models/
│   └── services/
├── stock-service/              # Go Stock servisi
│   ├── main.go
│   ├── handlers/
│   ├── models/
│   └── services/
└── shared/                     # Paylaşılan event modelleri
    └── events/
```

## 🎯 Özellikler

- ✅ Saga Choreography Pattern
- ✅ Inbox/Outbox Pattern (Event delivery garantisi)
- ✅ Asenkron event iletişimi
- ✅ Her servis kendi veritabanı
- ✅ Docker Compose ile kolay deployment
- ✅ JSON event formatı
- ✅ Sade ve takip edilebilir loglar

## 🔍 Event Türleri

### Order Events
- `OrderCreated`
- `OrderConfirmed`
- `OrderCancelled`

### Payment Events
- `PaymentRequested`
- `PaymentProcessed`
- `PaymentFailed`

### Stock Events
- `StockReserved`
- `StockReservationFailed`
