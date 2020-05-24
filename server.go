package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type GLInfo struct {
	bus     int
	address int
}

type SubscribeInfo struct {
	sessionID   int
	infoChannel chan interface{}
}

type UnsubscribeInfo struct {
	sessionID int
}

var subscriptions sync.Map

func runTcpServer(port int) {
	serverSocket, err := net.Listen("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Unable to open open listening socket (%s)", err.Error())
	}
	defer serverSocket.Close()

	subscriptionChannel := make(chan interface{})
	//	infoChannel := make(chan interface{})

	go handleRegisterListener(subscriptionChannel)

	go broadcast()

	for {
		connection, err := serverSocket.Accept()
		if err != nil {
			log.Fatalf("Error accepting incoming socket connection (%s)", err.Error())
		}
		NewTcpConnector(connection, subscriptionChannel).Start()
	}
}

func handleRegisterListener(subscriptionChannnel chan interface{}) {
	for {
		switch request := (<-subscriptionChannnel).(type) {
		case SubscribeInfo:
			subscriptions.Store(request.sessionID, request.infoChannel)
			log.Printf("Subscription added for session %d", request.sessionID)
		case UnsubscribeInfo:
			subscriptions.Delete(request.sessionID)
			log.Printf("Subscription removed for session %d", request.sessionID)
		}
	}
}

func broadcast() {
	duration, _ := time.ParseDuration("1s")
	for {
		time.Sleep(duration)
		subscriptions.Range(func(key, value interface{}) bool {
			subscriptor := value.(chan interface{})
			subscriptor <- "INFO 1 2 3 4"
			return true
		})
	}
}