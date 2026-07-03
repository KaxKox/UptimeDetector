package main

import (
	"fmt"
	"log"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://broker.emqx.io:1883")
	opts.SetClientID("gocheck_student_123")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Błąd połączenia MQTT: ", token.Error())
	}
	fmt.Println("Połączono z publicznym brokerem MQTT!")

	topic := "uczelnia/gocheck/alerts"
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Odebrano wiadomość MQTT: %s\n", msg.Payload())
	})

	tekst := "Zaczynamy monitorowanie serwerów! Awaria zlikwidowana."
	client.Publish(topic, 0, false, tekst)
	fmt.Println("Wysłano wiadomość na kanał:", topic)

	time.Sleep(2 * time.Second)
	client.Disconnect(250)
}