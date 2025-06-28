using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace OrderService.Models;

public class Order
{
    public string Id { get; set; } = Guid.NewGuid().ToString();
    public string CustomerId { get; set; } = string.Empty;
    public decimal TotalAmount { get; set; }
    public string Status { get; set; } = "Created";
    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    public DateTime? UpdatedAt { get; set; }
    
    public List<OrderItem> Items { get; set; } = new();
}

public class OrderItem
{
    public string Id { get; set; } = Guid.NewGuid().ToString();
    public string ProductId { get; set; } = string.Empty;
    public int Quantity { get; set; }
    public decimal Price { get; set; }
    public string OrderId { get; set; } = string.Empty;
    [JsonIgnore]
    public Order Order { get; set; } = null!;
}

public class CreateOrderRequest
{
    [Required]
    public string CustomerId { get; set; } = string.Empty;
    
    [Required]
    public List<CreateOrderItemRequest> Items { get; set; } = new();
}

public class CreateOrderItemRequest
{
    [Required]
    public string ProductId { get; set; } = string.Empty;
    
    [Range(1, int.MaxValue)]
    public int Quantity { get; set; }
    
    [Range(0.01, double.MaxValue)]
    public decimal Price { get; set; }
} 