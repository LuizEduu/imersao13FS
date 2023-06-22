package kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type Producer struct {
	ConfigMap *ckafka.ConfigMap
}

func NewKafkaProducer(configMap *ckafka.ConfigMap) *Producer {
	return &Producer{
		ConfigMap: configMap,
	}
}

func (p *Producer) Publish(message any, key []byte, topic string) error {
	producer, err := ckafka.NewProducer(p.ConfigMap)

	if err != nil {
		return err
	}

	defer producer.Close()

	kafkaMessage := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic:     &topic,
			Partition: ckafka.PartitionAny,
		},
		Key:   key,
		Value: message.([]byte),
	}

	err = producer.Produce(kafkaMessage, nil)

	if err != nil {
		return err
	}

	return nil
}
