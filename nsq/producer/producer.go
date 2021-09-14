package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

type Message struct {
	Name      string
	Timestamp string
}

func main() {
	//The only valid way to instantiate the Config
	config := nsq.NewConfig()

	//Creating the Producer using NSQD Address
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}

	//Init topic name and message
	topic := "logs"

	msg := Message{
		Name:      "test",
		Timestamp: time.Now().String(),
	}

	//Convert message as []byte
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}

	//Publish the Message
	err = producer.Publish(topic, payload)
	if err != nil {
		log.Println(err)
	}
}
