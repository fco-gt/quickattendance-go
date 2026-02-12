package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"quickattendance-go/internal/config"
	"quickattendance-go/pkg/logger"
	"syscall"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}
	cfg := config.Load()
	logger.Setup(cfg.Env)

	conn, err := amqp.Dial(cfg.RabbitURL)
	if err != nil {
		slog.Error("Error conectando a RabbitMQ", "error", err)
		os.Exit(1) // Detenemos si no hay conexión
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		slog.Error("Error abriendo canal", "error", err)
		os.Exit(1)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("email_queue", true, false, false, false, nil)
	if err != nil {
		slog.Error("Error declarando cola", "error", err)
		os.Exit(1)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		slog.Error("Error registrando consumidor", "error", err)
		os.Exit(1)
	}

	go func() {
		for d := range msgs {
			var emailData map[string]string
			if err := json.Unmarshal(d.Body, &emailData); err != nil {
				slog.Error("Error decodificando mensaje", "error", err)
				continue
			}

			// LOG ESTRUCTURADO: Fácil de leer y de procesar por máquinas
			slog.Info("Enviando email",
				"to", emailData["to"],
				"subject", emailData["subject"],
				"body", emailData["body"],
			)

			// Aquí iría tu lógica real de SMTP
			slog.Info("Email enviado exitosamente", "to", emailData["to"])
		}
	}()

	slog.Info("Worker operativo", "queue", q.Name)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("Deteniendo worker...")
}
