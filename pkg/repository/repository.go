package repository

import (
	"L0test/model"
	"database/sql"
	"fmt"
	"log"
)

type Repository interface {
	OrdersFromDb()
	FindById(orderUID string) model.Order
	InsertOrder(order *model.Order) error
}

type orderRepository struct {
	All map[string]model.Order
}

func New() *orderRepository {
	return &orderRepository{
		All: make(map[string]model.Order),
	}
}

func (r *orderRepository) OrdersFromDb() {
	connStr := "user=dev password=dev dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rowsOrder, err := db.Query("SELECT * FROM orders")
	if err != nil {
		log.Fatal(err)
	}

	for rowsOrder.Next() {
		order := model.Order{}
		rowsOrder.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)

		rowsPayment, err := db.Query("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid=$1", &order.OrderUID)
		if err != nil {
			log.Fatalf("Error query rows payment: %s", err)
		}

		payment := model.Payment{}
		for rowsPayment.Next() {
			if err = rowsPayment.Scan(&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee); err != nil {
				log.Fatalf("Error scan rows payment: %s", err)
			}
		}

		rowsDelivery, err := db.Query("SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid=$1", &order.OrderUID)
		if err != nil {
			log.Fatalf("Error query rows delivery: %s", err)
		}

		delivery := model.Delivery{}
		for rowsDelivery.Next() {
			err = rowsDelivery.Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
			if err != nil {
				log.Fatalf("Error scan rows delivery: %s", err)
			}
		}

		rowsItems, err := db.Query("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid=$1", &order.OrderUID)
		if err != nil {
			log.Fatalf("Error query rows items: %s", err)
		}

		var items []model.Item
		for rowsItems.Next() {
			item := model.Item{}
			err = rowsItems.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
			if err != nil {
				fmt.Printf("Error scan items: %s", err)
				continue
			}

			items = append(items, item)
		}

		o := model.Order{
			OrderUID:          order.OrderUID,
			TrackNumber:       order.TrackNumber,
			Entry:             order.Entry,
			Delivery:          delivery,
			Payment:           payment,
			Items:             items,
			Locale:            order.Locale,
			InternalSignature: order.InternalSignature,
			CustomerID:        order.CustomerID,
			DeliveryService:   order.DeliveryService,
			ShardKey:          order.ShardKey,
			SmID:              order.SmID,
			DateCreated:       order.DateCreated,
			OofShard:          order.OofShard,
		}

		r.All[order.OrderUID] = o
	}
}

func (r *orderRepository) FindById(orderUID string) model.Order {
	order, ok := r.All[orderUID]
	if !ok {
		log.Println("Not this order: " + orderUID)
		return model.Order{}
	}
	return order
}

func (r *orderRepository) InsertOrder(order *model.Order) error {
	connStr := "user=dev password=dev dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connect to DB: %s", err)
		return err
	}

	defer db.Close()

	db.QueryRow(`INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2,$3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)

	db.QueryRow(`INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3,$4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)

	db.QueryRow(`INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3,$4, $5, $6, $7, $8)`,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)

	for _, item := range order.Items {
		db.QueryRow(`INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3,$4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
	}
	r.All[order.OrderUID] = *order

	return nil
}
