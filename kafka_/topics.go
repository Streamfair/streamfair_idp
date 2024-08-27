package kafka_

import (
    "context"
    "fmt"

    "github.com/confluentinc/confluent-kafka-go/kafka"
)

// CreateTopics creates Kafka topics using the admin client.
func CreateTopics(ctx context.Context, broker string, topics ...string) error {
    adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
        "bootstrap.servers": broker,
    })
    if err != nil {
        return fmt.Errorf("failed to create admin client: %v", err)
    }
    defer adminClient.Close()

    var topicSpecs []kafka.TopicSpecification
    for _, topic := range topics {
        topicSpecs = append(topicSpecs, kafka.TopicSpecification{
            Topic:             topic,
            NumPartitions:     3, // Adjust as needed
            ReplicationFactor: 1, // Adjust as needed
        })
    }

    results, err := adminClient.CreateTopics(ctx, topicSpecs, nil)
    if err != nil {
        return fmt.Errorf("failed to create topics: %v", err)
    }

    for _, result := range results {
        if result.Error.Code() != kafka.ErrNoError {
            return fmt.Errorf("failed to create topic %s: %v", result.Topic, result.Error)
        }
        fmt.Printf("Topic %s created successfully\n", result.Topic)
    }

    return nil
}
