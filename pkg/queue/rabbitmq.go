package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/prefeitura-rio/app-notification-core/internal/config"
	"github.com/prefeitura-rio/app-notification-core/internal/entity"
)

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *config.Config
}

type NotificationMessage struct {
	Notification *entity.Notification `json:"notification"`
	Timestamp    time.Time             `json:"timestamp"`
	RetryCount   int                   `json:"retry_count"`
}

func NewRabbitMQClient(cfg *config.Config) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declarar fila com configurações de durabilidade
	_, err = channel.QueueDeclare(
		cfg.RabbitMQ.QueueNotifications, // name
		true,                              // durable
		false,                             // delete when unused
		false,                             // exclusive
		false,                             // no-wait
		amqp.Table{
			"x-message-ttl":             int32(3600000), // 1 hora
			"x-max-length":              int32(100000),  // máximo 100k mensagens
			"x-dead-letter-exchange":    "notifications.dlx",
			"x-dead-letter-routing-key": "notifications.dlq",
		},
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Declarar exchange para Dead Letter Queue
	err = channel.ExchangeDeclare(
		"notifications.dlx", // name
		"direct",            // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		log.Printf("Warning: Failed to declare DLX: %v", err)
	}

	// Declarar Dead Letter Queue
	_, err = channel.QueueDeclare(
		"notifications.dlq", // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,
	)
	if err != nil {
		log.Printf("Warning: Failed to declare DLQ: %v", err)
	}

	// Bind DLQ ao exchange
	err = channel.QueueBind(
		"notifications.dlq", // queue name
		"notifications.dlq", // routing key
		"notifications.dlx", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Printf("Warning: Failed to bind DLQ: %v", err)
	}

	// Configurar QoS (prefetch)
	err = channel.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Printf("Warning: Failed to set QoS: %v", err)
	}

	log.Printf("✅ RabbitMQ connected to %s", cfg.RabbitMQ.URL)
	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
		config:  cfg,
	}, nil
}

// PublishNotification publica uma notificação na fila
func (r *RabbitMQClient) PublishNotification(notification *entity.Notification) error {
	message := NotificationMessage{
		Notification: notification,
		Timestamp:    time.Now(),
		RetryCount:   0,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.channel.PublishWithContext(
		ctx,
		"",                                  // exchange
		r.config.RabbitMQ.QueueNotifications, // routing key
		false,                               // mandatory
		false,                               // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
			MessageId:    notification.ID.String(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("📤 Notification %s published to queue", notification.ID)
	return nil
}

// ConsumeNotifications consome mensagens da fila
func (r *RabbitMQClient) ConsumeNotifications(handler func(*NotificationMessage) error) error {
	msgs, err := r.channel.Consume(
		r.config.RabbitMQ.QueueNotifications, // queue
		"",                                     // consumer
		false,                                  // auto-ack
		false,                                  // exclusive
		false,                                  // no-local
		false,                                  // no-wait
		nil,                                    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("🔄 Consumer started, waiting for messages...")

	for msg := range msgs {
		var notifMsg NotificationMessage
		if err := json.Unmarshal(msg.Body, &notifMsg); err != nil {
			log.Printf("❌ Failed to unmarshal message: %v", err)
			msg.Nack(false, false) // Envia para DLQ
			continue
		}

		log.Printf("📥 Processing notification %s (retry: %d)", notifMsg.Notification.ID, notifMsg.RetryCount)

		// Processar mensagem
		if err := handler(&notifMsg); err != nil {
			log.Printf("❌ Failed to process notification %s: %v", notifMsg.Notification.ID, err)

			// Retry logic
			if notifMsg.RetryCount < 3 {
				// Republicar com retry incrementado
				notifMsg.RetryCount++
				body, _ := json.Marshal(notifMsg)

				r.channel.Publish(
					"",
					r.config.RabbitMQ.QueueNotifications,
					false,
					false,
					amqp.Publishing{
						DeliveryMode: amqp.Persistent,
						ContentType:  "application/json",
						Body:         body,
					},
				)
				msg.Ack(false)
				log.Printf("🔄 Notification %s requeued (retry %d/3)", notifMsg.Notification.ID, notifMsg.RetryCount)
			} else {
				// Após 3 tentativas, envia para DLQ
				msg.Nack(false, false)
				log.Printf("💀 Notification %s sent to DLQ after 3 retries", notifMsg.Notification.ID)
			}
			continue
		}

		// Sucesso
		msg.Ack(false)
		log.Printf("✅ Notification %s processed successfully", notifMsg.Notification.ID)
	}

	return nil
}

// GetQueueStats retorna estatísticas da fila
func (r *RabbitMQClient) GetQueueStats() (map[string]interface{}, error) {
	queue, err := r.channel.QueueDeclarePassive(
		r.config.RabbitMQ.QueueNotifications,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue stats: %w", err)
	}

	dlq, _ := r.channel.QueueDeclarePassive(
		"notifications.dlq",
		true,
		false,
		false,
		false,
		nil,
	)

	stats := map[string]interface{}{
		"queue_name":    queue.Name,
		"messages":      queue.Messages,
		"consumers":     queue.Consumers,
		"dlq_messages":  dlq.Messages,
		"last_checked":  time.Now(),
	}

	return stats, nil
}

// PurgeQueue limpa todas as mensagens da fila
func (r *RabbitMQClient) PurgeQueue() error {
	_, err := r.channel.QueuePurge(r.config.RabbitMQ.QueueNotifications, false)
	if err != nil {
		return fmt.Errorf("failed to purge queue: %w", err)
	}
	log.Printf("🗑️ Queue %s purged", r.config.RabbitMQ.QueueNotifications)
	return nil
}

// Close fecha a conexão com RabbitMQ
func (r *RabbitMQClient) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
