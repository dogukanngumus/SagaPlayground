using Newtonsoft.Json;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using System.Text;

namespace OrderService.Services;

public class EventConsumerService : BackgroundService
{
    private readonly IConnection _connection;
    private readonly IModel _channel;
    private readonly OrderBusinessService _orderService;
    private readonly ILogger<EventConsumerService> _logger;

    public EventConsumerService(OrderBusinessService orderService, ILogger<EventConsumerService> logger)
    {
        _orderService = orderService;
        _logger = logger;

        var factory = new ConnectionFactory
        {
            HostName = Environment.GetEnvironmentVariable("RABBITMQ_HOST") ?? "localhost",
            UserName = Environment.GetEnvironmentVariable("RABBITMQ_USER") ?? "admin",
            Password = Environment.GetEnvironmentVariable("RABBITMQ_PASS") ?? "admin123"
        };

        _connection = factory.CreateConnection();
        _channel = _connection.CreateModel();

        // Declare queue for order events
        _channel.QueueDeclare("order_events", durable: true, exclusive: false, autoDelete: false);

        // Bind queue to exchange with correct routing keys
        _channel.QueueBind("order_events", "saga.events", "stock.reserved");
        _channel.QueueBind("order_events", "saga.events", "stock.reservation.failed");
    }

    protected override Task ExecuteAsync(CancellationToken stoppingToken)
    {
        var consumer = new EventingBasicConsumer(_channel);
        consumer.Received += async (model, ea) =>
        {
            var body = ea.Body.ToArray();
            var message = Encoding.UTF8.GetString(body);
            var routingKey = ea.RoutingKey;

            _logger.LogInformation("Received message: {Message} with routing key: {RoutingKey}", message, routingKey);

            try
            {
                var eventData = JsonConvert.DeserializeObject<Dictionary<string, object>>(message);
                
                if (eventData != null && eventData.ContainsKey("eventType"))
                {
                    var eventType = eventData["eventType"].ToString();
                    var orderId = eventData["orderId"].ToString();

                    switch (eventType)
                    {
                        case "StockReserved":
                            await _orderService.UpdateOrderStatusAsync(orderId, "Confirmed");
                            _logger.LogInformation("Order confirmed: {OrderId}", orderId);
                            break;

                        case "StockReservationFailed":
                            await _orderService.UpdateOrderStatusAsync(orderId, "Cancelled");
                            _logger.LogInformation("Order cancelled: {OrderId}", orderId);
                            break;
                    }
                }
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Error processing message: {Message}", message);
            }

            _channel.BasicAck(ea.DeliveryTag, false);
        };

        _channel.BasicConsume(queue: "order_events",
                             autoAck: false,
                             consumer: consumer);

        return Task.CompletedTask;
    }

    public override void Dispose()
    {
        _channel?.Dispose();
        _connection?.Dispose();
        base.Dispose();
    }
} 