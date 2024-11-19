package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"

	messagingqueue "github.com/gnanasuryateja/golib/datastore/messaging_queue"
)

type kafkaStore struct {
	client  sarama.Client
	lock    sync.Mutex
	brokers []string
}

// creates a new kafka client
func NewKafkaStoreClient(ctx context.Context, brokers []string) (messagingqueue.MessageQueue, error) {
	return kafkaStore{
		brokers: brokers,
	}, nil
}

// GetKafkaClient returns an existing Kafka client if available, or creates a new one
func (k *kafkaStore) GetKafkaClient(brokers []string) (sarama.Client, error) {
	k.lock.Lock()
	defer k.lock.Unlock()

	// check if client already exists and is healthy
	if k.client != nil && !k.client.Closed() {
		fmt.Println("Using existing Kafka client")
		return k.client, nil
	}

	// create a new kafka client if it doesn't exist or is closed
	fmt.Println("Creating a new Kafka client")
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to Kafka broker: %v", err)
	}

	k.client = client
	return k.client, nil
}

// checks the connection to kafka and return error if any
func (k kafkaStore) HealthCheck(ctx context.Context) error {

	// get the kafka client
	client, err := k.GetKafkaClient(k.brokers)
	if err != nil {
		return err
	}

	// refresh metadata to ensure Kafka is reachable
	err = client.RefreshMetadata()
	if err != nil {
		return fmt.Errorf("failed to refresh Kafka metadata: %v", err)
	}

	// check if we can retrieve broker information
	brokersList := client.Brokers()
	if len(brokersList) == 0 {
		return fmt.Errorf("no brokers available")
	}

	fmt.Println("Kafka is healthy, brokers:", brokersList)
	return nil
}

// sends a message to a topic
func (k kafkaStore) ProduceMessage(ctx context.Context, args ...any) error {
	// Get or create a Kafka client
	client, err := k.GetKafkaClient(k.brokers)
	if err != nil {
		return err
	}

	// Create a new Kafka sync producer using the existing client
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// get the topic and message
	topic := args[0].(string)
	message := args[1].(string)

	// Prepare the message to send
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	// Send the message
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	return nil
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Printf("Message claimed: value = %s, timestamp = %v, topic = %s\n",
			string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}
	return nil
}

// receives a message from a topic
func (k kafkaStore) ConsumeMessage(ctx context.Context, args ...any) (any, error) {
	// Get or create a Kafka client
	client, err := k.GetKafkaClient(k.brokers)
	if err != nil {
		return nil, err
	}

	// get the topic and groupId
	topic := args[0].(string)
	groupId := args[1].(string)

	// Create a new Kafka consumer group
	consumerGroup, err := sarama.NewConsumerGroupFromClient(groupId, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// Create a consumer handler
	consumer := Consumer{
		ready: make(chan bool),
	}

	// Consume messages in a loop
	for {
		err := consumerGroup.Consume(ctx, []string{topic}, &consumer)
		if err != nil {
			return nil, fmt.Errorf("error while consuming messages: %v", err)
		}

		// Check if the consumer is ready
		if !<-consumer.ready {
			return nil, fmt.Errorf("consumer failed to be ready")
		}
	}
}
