package services

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type MessageBrokerService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewMessageBrokerService() (*MessageBrokerService, error) {
	host := getEnv("RABBITMQ_HOST", "localhost")
	user := getEnv("RABBITMQ_USER", "admin")
	password := getEnv("RABBITMQ_PASS", "admin123")

	url := "amqp://" + user + ":" + password + "@" + host + ":5672/"

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"saga.events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"payment_events",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		"order.created",
		"saga.events",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	log.Println("Message broker connected successfully")
	return &MessageBrokerService{
		conn:    conn,
		channel: ch,
	}, nil
}

func (m *MessageBrokerService) PublishMessage(exchange, routingKey string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return m.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (m *MessageBrokerService) ConsumeMessages(queueName string, handler func([]byte) error) error {
	msgs, err := m.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
			if err := handler(d.Body); err != nil {
				log.Printf("Error processing message: %v", err)
			}
		}
	}()

	return nil
}

func (m *MessageBrokerService) Close() {
	if m.channel != nil {
		m.channel.Close()
	}
	if m.conn != nil {
		m.conn.Close()
	}
}
