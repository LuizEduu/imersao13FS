package kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type Consumer struct {
	ConfigMap *ckafka.ConfigMap
	Topics    []string
}

func NewConsumer(configMap *ckafka.ConfigMap, topics []string) *Consumer {
	return &Consumer{
		ConfigMap: configMap,
		Topics:    topics,
	}
}

func (c *Consumer) Consume(messageChanel chan *ckafka.Message) error {
	consumer, err := ckafka.NewConsumer(c.ConfigMap)

	if err != nil {
		panic(err)
	}

	err = consumer.SubscribeTopics(c.Topics, nil)

	if err != nil {
		panic(err)
	}

	for {
		message, err := consumer.ReadMessage(-1)

		if err == nil {
			messageChanel <- message
		}
	}

}
