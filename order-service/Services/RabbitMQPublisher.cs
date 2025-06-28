using Microsoft.EntityFrameworkCore;
using Newtonsoft.Json;
using OrderService.Data;
using OrderService.Models;
using RabbitMQ.Client;
using System.Text;

namespace OrderService.Services;

public class RabbitMQPublisher : IMessagePublisher, IDisposable
{
    private readonly IConnection _connection;
    private readonly IModel _channel;
    private readonly OrderDbContext _context;
    private readonly ILogger<RabbitMQPublisher> _logger;

    public RabbitMQPublisher(OrderDbContext context, ILogger<RabbitMQPublisher> logger)
    {
        _context = context;
        _logger = logger;

        var factory = new ConnectionFactory
        {
            HostName = Environment.GetEnvironmentVariable("RABBITMQ_HOST") ?? "localhost",
            UserName = Environment.GetEnvironmentVariable("RABBITMQ_USER") ?? "admin",
            Password = Environment.GetEnvironmentVariable("RABBITMQ_PASS") ?? "admin123"
        };

        _connection = factory.CreateConnection();
        _channel = _connection.CreateModel();

        _channel.ExchangeDeclare("saga.events", ExchangeType.Topic, durable: true);
        
        _logger.LogInformation("RabbitMQ publisher initialized");
    }

    public async Task PublishAsync<T>(T message, string exchange, string routingKey) where T : class
    {
        try
        {
            var outboxMessage = new OutboxMessage
            {
                EventType = typeof(T).Name,
                EventData = JsonConvert.SerializeObject(message),
                Exchange = exchange,
                RoutingKey = routingKey
            };

            _context.OutboxMessages.Add(outboxMessage);
            await _context.SaveChangesAsync();

            _logger.LogInformation("Message saved to outbox: {EventType}", typeof(T).Name);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error saving message to outbox");
            throw;
        }
    }

    public async Task PublishOutboxAsync()
    {
        try
        {
            var unprocessedMessages = await _context.OutboxMessages
                .Where(m => !m.IsProcessed)
                .OrderBy(m => m.CreatedAt)
                .ToListAsync();

            foreach (var message in unprocessedMessages)
            {
                try
                {
                    var body = Encoding.UTF8.GetBytes(message.EventData);
                    _channel.BasicPublish(
                        exchange: message.Exchange,
                        routingKey: message.RoutingKey,
                        basicProperties: null,
                        body: body);

                    message.IsProcessed = true;
                    message.ProcessedAt = DateTime.UtcNow;

                    _logger.LogInformation("Published outbox message: {EventType}", message.EventType);
                }
                catch (Exception ex)
                {
                    _logger.LogError(ex, "Error publishing outbox message: {MessageId}", message.Id);
                }
            }

            await _context.SaveChangesAsync();
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error processing outbox messages");
        }
    }

    public void Dispose()
    {
        _channel?.Dispose();
        _connection?.Dispose();
    }
} 