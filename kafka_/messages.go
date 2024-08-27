package kafka_

// import (
// 	"fmt"
// 	"log"

// 	"google.golang.org/protobuf/proto"
// )

// // handleKafkaMessage routes Kafka messages to appropriate handlers
// func handleKafkaMessage(key, value []byte) {
// 	handler, ok := messageHandlers[string(key)]
// 	if !ok {
// 		log.Printf("No handler found for message with key: %s", key)
// 		return
// 	}

// 	if err := handler(key, value); err != nil {
// 		log.Printf("Error handling message with key %s: %v", key, err)
// 	}
// }

// // Define a handler function type
// type MessageHandler func(key, value []byte) error

// // Map of handlers for different message types
// var messageHandlers = map[string]MessageHandler{
// 	"user_updated":         handleUserUpdated,
// 	"account_updated":      handleAccountUpdated,
// 	"account_type_updated": handleAccountTypeUpdated,
// 	"user_deleted":         handleUserDeleted,
// 	"account_deleted":      handleAccountDeleted,
// 	"account_type_deleted": handleAccountTypeDeleted,
// 	// Add more handlers for different message types here
// }

// // handleAccountCreated is an example handler for the "account_created" message type
// func handleAccountCreated(key, value []byte) error {
// 	// Example: Unmarshal the protobuf message
// 	accountCreated := &YourProtoMessage{}
// 	err := proto.Unmarshal(value, accountCreated)
// 	if err != nil {
// 		return fmt.Errorf("failed to unmarshal account_created message: %v", err)
// 	}

// 	// Example: Process the message (e.g., save to database)
// 	fmt.Printf("Received account_created message. Account ID: %s\n", key)
// 	// Your processing logic here

// 	return nil
// }

// // handleUserCreated is an example handler for the "user_created" message type
// func handleUserCreated(key, value []byte) error {
// 	// Example: Unmarshal the protobuf message
// 	userCreated := &YourProtoMessage{}
// 	err := proto.Unmarshal(value, userCreated)
// 	if err != nil {
// 		return fmt.Errorf("failed to unmarshal user_created message: %v", err)
// 	}

// 	// Example: Process the message (e.g., send welcome email)
// 	fmt.Printf("Received user_created message. User ID: %s\n", key)
// 	// Your processing logic here

// 	return nil
// }

// // handleAccountTypeCreated is an example handler for the "account_type_created" message type
// func handleAccountTypeCreated(key, value []byte) error {
// 	// Example: Unmarshal the protobuf message
// 	accountTypeCreated := &YourProtoMessage{}
// 	err := proto.Unmarshal(value, accountTypeCreated)
// 	if err != nil {
// 		return fmt.Errorf("failed to unmarshal account_type_created message: %v", err)
// 	}

// 	// Example: Process the message (e.g., update cache)
// 	fmt.Printf("Received account_type_created message. Account Type ID: %s\n", key)
// 	// Your processing logic here

// 	return nil
// }
