package kafka

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/your-org/error-simulator/config"
	"github.com/your-org/error-simulator/models"
)

var (
	producer *kafka.Producer
	once     sync.Once
	enabled  bool
)

// InitProducer creates the Kafka producer from config. Call from main before starting the server.
// If bootstrap servers are empty or init fails, publishing is disabled and errors are only logged.
func InitProducer(cfg *config.Config) {
	once.Do(func() {
		if cfg.KafkaBootstrapServers == "" {
			log.Println("[kafka] KAFKA_BOOTSTRAP_SERVERS not set; skipping Kafka producer (events will not be published)")
			enabled = false
			return
		}
		p, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": cfg.KafkaBootstrapServers,
		})
		if err != nil {
			log.Printf("[kafka] failed to create producer: %v; events will not be published", err)
			enabled = false
			return
		}
		producer = p
		enabled = true
		log.Printf("[kafka] producer connected to %s, topic %s", cfg.KafkaBootstrapServers, cfg.KafkaTopic)
	})
}

// PublishErrorEvent serializes the event to JSON and publishes to the configured topic.
// Never panics; if Kafka is unreachable or not configured, logs and returns an error.
func PublishErrorEvent(cfg *config.Config, event models.ErrorEvent) error {
	if !enabled || producer == nil {
		log.Println("[kafka] producer not available; skipping publish")
		return nil
	}
	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("[kafka] marshal error: %v", err)
		return err
	}
	topic := cfg.KafkaTopic
	if topic == "" {
		topic = "service.errors"
	}
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
	}, nil)
	if err != nil {
		log.Printf("[kafka] produce error: %v", err)
		return err
	}
	log.Printf("[kafka] event published to %s", topic)
	return nil
}
