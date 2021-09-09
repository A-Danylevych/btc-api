package main

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type MessageBody struct {
	Level   string `json:"level"`
	Message string `json:"msg"`
	Time    string `json:"time"`
}

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	if err != nil {
		fmt.Println(err)
	}

	messages, err := ch.Consume(
		"BTC-API",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Print(err.Error())
	}

	forever := make(chan bool)
	go func() {
		for mes := range messages {
			message := &MessageBody{}
			json.Unmarshal(mes.Body, &message)
			if message.Level == "error" {
				fmt.Printf("Message: %s\nTime: %s\n", message.Message, message.Time)

			}
		}
	}()

	<-forever
}
