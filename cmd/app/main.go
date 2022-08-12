package main

import (
	"L0test/model"
	"L0test/pkg/handler"
	"L0test/pkg/repository"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
)

func main() {
	orderRep := repository.New()
	orderRep.OrdersFromDb()
	respOrder := handler.RegOrder{Order: orderRep}

	go listenNats(&respOrder)
	router := httprouter.New()
	router.GET("/orders/:id", respOrder.Response)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func listenNats(o *handler.RegOrder) {
	// Connect to a server
	nc, err := nats.Connect("nats://0.0.0.0:4222")
	if err != nil {
		log.Fatalf("Error to connect nats: %s", err)
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatalf("Error encode: %s", err)
	}

	// Simple Async Subscriber
	ec.Subscribe("orders", func(msg *nats.Msg) {
		var order model.Order
		// Unmarshal JSON that represents the Order data
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Fatalf("Error unmarshal data: %s", err)
			return
		}
		log.Println("Received order: ", order.OrderUID)

		// Insert Order into DB
		err := o.Order.InsertOrder(&order)
		if err != nil {
			log.Fatalf("Error to insert order: %s", err)
			return
		}

		log.Println("Added new order: ", order.OrderUID)
	})
}
