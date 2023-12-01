package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

type Goods struct {
	Text string `json:"text"`
	Num  int    `json:"num"`
}

func msgCallback(msg *nats.Msg) {
	// Десериализация полученных байтов
	var receivedMessage Goods
	buf := bytes.NewBuffer(msg.Data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&receivedMessage)
	if err != nil {
		log.Println("Error decoding message:", err)
		return
	}

	fmt.Println(receivedMessage)
}

const (
	subject = "test-subject"
)

func main() {
	options := &server.Options{}
	ns, err := server.NewServer(options)

	if err != nil {
		log.Fatal(err)
	}

	go ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		log.Fatal("Not ready")
	}

	fmt.Println("Server started")

	nc, err := nats.Connect(ns.ClientURL())

	if err != nil {
		log.Fatal(err)
	}

	nc.Subscribe(subject, msgCallback)

	// Создаем сообщение и сериализуем его в поток байт с использованием gob
	message := Goods{Text: "Hi!", Num: 1}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(message)
	if err != nil {
		log.Fatal("Error encoding message:", err)
	}

	// Отправляем сериализованное сообщение
	nc.Publish(subject, buf.Bytes())

	ns.Shutdown()
	ns.WaitForShutdown()
}
