package kafka_

import (
    "fmt"

    "github.com/confluentinc/confluent-kafka-go/kafka"
    "github.com/rs/zerolog/log"
)

// PublishMessage publishes a message to Kafka.
func PublishMessage(broker string, topic string, key []byte, value []byte) error {
    p, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": broker,
    })
    if err != nil {
        return fmt.Errorf("failed to create Kafka producer: %v", err)
    }
    defer p.Close()

    deliveryChan := make(chan kafka.Event, 1)
    err = p.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Key:            key,
        Value:          value,
    }, deliveryChan)

    e := <-deliveryChan
    msg := e.(*kafka.Message)
    if msg.TopicPartition.Error != nil {
        return fmt.Errorf("failed to deliver message: %v", msg.TopicPartition.Error)
    }

    log.Info().Msgf("Message published: %s", string(value))
    return nil
}