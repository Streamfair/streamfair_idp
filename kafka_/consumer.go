package kafka_

// import (
// 	"os"

// 	"github.com/confluentinc/confluent-kafka-go/kafka"
// 	"github.com/rs/zerolog/log"
// )

// func SetupKafkaConsumer() *kafka.Consumer {
// 	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
// 		"bootstrap.servers":  "localhost:9092", // Update with your Kafka broker(s)
// 		"group.id":           "my-consumer-group",
// 		"auto.offset.reset":  "earliest",
// 		"enable.auto.commit": "false",
// 	})
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to create Kafka consumer")
// 	}
// 	err = consumer.SubscribeTopics([]string{"account_created"}, nil)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("failed to subscribe to Kafka topics")
// 	}
// 	return consumer
// }

// func ConsumeKafkaMessages(consumer *kafka.Consumer, sigchan chan os.Signal) {
// 	for {
// 		select {
// 		case sig := <-sigchan:
// 			log.Info().Msgf("Caught signal %v: terminating", sig)
// 			return
// 		default:
// 			ev := consumer.Poll(100)
// 			if ev == nil {
// 				continue
// 			}
// 			switch e := ev.(type) {
// 			case *kafka.Message:
// 				log.Info().Msgf("Received message from Kafka: Key: %s, Value: %s", string(e.Key), string(e.Value))
// 				handleKafkaMessage(e.Key, e.Value) // Route message to appropriate handler
// 			case kafka.Error:
// 				log.Error().Msgf("Kafka consumer error: %v", e)
// 			default:
// 				log.Warn().Msgf("Ignoring event: %v", e)
// 			}
// 		}
// 	}
// }
