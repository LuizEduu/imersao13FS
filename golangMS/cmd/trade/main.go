package main

import (
	"encoding/json"
	"fmt"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/luizeduu/imersao13/go/infra/kafka"
	"github.com/luizeduu/imersao13/go/internal/dto"
	"github.com/luizeduu/imersao13/go/internal/market/entities"
	"github.com/luizeduu/imersao13/go/internal/market/transformers"
)

func main() {
	ordersIn := make(chan *entities.Order)
	ordersOut := make(chan *entities.Order)
	waitGroup := &sync.WaitGroup{}

	defer waitGroup.Wait()

	kafkaMessageChannel := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "tradeGroup",
		"auto.offset.reset": "latest",
	}

	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"orders"})
	go kafka.Consume(kafkaMessageChannel)

	book := entities.NewBook(ordersIn, ordersOut, waitGroup)
	go book.Trade()

	go func() {
		for message := range kafkaMessageChannel {
			waitGroup.Add(1)
			fmt.Println(string(message.Value))
			tradeInput := dto.TradeInput{}

			err := json.Unmarshal(message.Value, &tradeInput)

			if err != nil {
				panic(err)
			}

			order := transformers.TransformInput(tradeInput)

			ordersIn <- order
		}
	}()

	for res := range ordersOut {
		output := transformers.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", "   ")

		fmt.Println(string(outputJson))
		if err != nil {
			fmt.Println(err)
		}

		producer.Publish(outputJson, []byte("transactions"), "transactions")

	}

}
