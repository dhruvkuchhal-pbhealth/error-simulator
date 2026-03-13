package kafka

import (
	"encoding/json"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/your-org/error-simulator/config"
	"github.com/your-org/error-simulator/logger"
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
			logger.Log.Warn().Msg("KAFKA_BOOTSTRAP_SERVERS not set; skipping Kafka producer (events will not be published)")
			enabled = false
			return
		}
		p, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": cfg.KafkaBootstrapServers,
		})
		if err != nil {
			logger.Log.Error().Err(err).Msg("kafka producer init failed; events will not be published")
			enabled = false
			return
		}
		producer = p
		enabled = true
		logger.Log.Info().
			Str("bootstrap", cfg.KafkaBootstrapServers).
			Str("topic", cfg.KafkaTopic).
			Msg("kafka producer connected")
	})
}

// PublishErrorEvent serializes the event to JSON and publishes to the configured topic.
// Never panics; if Kafka is unreachable or not configured, logs and returns an error.
func PublishErrorEvent(cfg *config.Config, event models.ErrorEvent) error {
	if !enabled || producer == nil {
		logger.Log.Warn().Msg("kafka producer not available; skipping publish")
		return nil
	}
	payload, err := json.Marshal(event)
	if err != nil {
		logger.Log.Error().Err(err).Msg("kafka event marshal failed")
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
		logger.Log.Error().Err(err).Str("topic", topic).Msg("kafka produce failed")
		return err
	}
	logger.Log.Info().Str("topic", topic).Str("error_message", event.ErrorMessage).Msg("error event published to kafka")
	return nil
}
