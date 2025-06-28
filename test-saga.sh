#!/bin/bash

echo "🚀 Saga Choreography Demo Test Script"
echo "======================================"

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 10

# Test 1: Create an order
echo ""
echo "📦 Test 1: Creating an order..."
ORDER_RESPONSE=$(curl -s -X POST http://localhost:5100/api/orders \
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
  }')

echo "Order Response: $ORDER_RESPONSE"

# Extract order ID from response
ORDER_ID=$(echo $ORDER_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

if [ -z "$ORDER_ID" ]; then
    echo "❌ Failed to create order"
    exit 1
fi

echo "✅ Order created with ID: $ORDER_ID"

# Test 2: Check order status
echo ""
echo "🔍 Test 2: Checking order status..."
sleep 5

ORDER_STATUS=$(curl -s http://localhost:5100/api/orders/$ORDER_ID)
echo "Order Status: $ORDER_STATUS"

# Test 3: Check payments
echo ""
echo "💳 Test 3: Checking payments..."
sleep 3

PAYMENTS=$(curl -s http://localhost:5001/api/payments)
echo "Payments: $PAYMENTS"

# Test 4: Check products
echo ""
echo "📦 Test 4: Checking products..."
sleep 3

PRODUCTS=$(curl -s http://localhost:5002/api/products)
echo "Products: $PRODUCTS"

echo ""
echo "🎉 Test completed! Check the logs to see the event flow:"
echo "   - Order Service: http://localhost:5100"
echo "   - Payment Service: http://localhost:5001"
echo "   - Stock Service: http://localhost:5002"
echo "   - RabbitMQ Management: http://localhost:15672 (admin/admin123)" 