package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewRabbitMQProducer(url string, queueName string) (*RabbitMQProducer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName, // nombre
		true,      // durable (sobrevive a reinicios de RabbitMQ)
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &RabbitMQProducer{
		conn:    conn,
		channel: ch,
		queue:   queueName,
	}, nil
}

// PublishEmail cumple con la interface domain.NotificationProvider
func (p *RabbitMQProducer) PublishEmail(ctx context.Context, to string, subject string, body string) error {
	message := map[string]string{
		"to":      to,
		"subject": subject,
		"body":    body,
	}

	// Convertir el mapa a JSON (bytes) para RabbitMQ
	bodyBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal email message: %w", err)
	}

	err = p.channel.PublishWithContext(ctx,
		"",      // exchange (default)
		p.queue, // routing key (nombre de la cola)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyBytes,
		})

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf(" [x] Mensaje enviado a la cola %s para: %s", p.queue, to)
	return nil
}

func (p *RabbitMQProducer) Close() {
	p.channel.Close()
	p.conn.Close()
}
