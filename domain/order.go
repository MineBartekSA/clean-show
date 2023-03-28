package domain

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Order struct {
	DBModel
	Status          OrderStatus           `db:"status"`
	OrderBy         uint                  `db:"order_by"`
	ShippingAddress string                `db:"shipping_address"`
	InvoiceAddress  string                `db:"invoice_address"`
	Products        DBArray[ProductOrder] `db:"products"`
	ShippingPrice   float64               `db:"shipping_price"`
	Total           float64               `db:"total"`
}

type OrderStatus int

const (
	OrderStatusCreated OrderStatus = iota + 1
	OrderStatusPaid
	OrderStatusInRealisation
	OrderStatusShipped
	OrderStatusCompleted
	OrderStatusCanceled
)

type ProductOrder struct {
	ProductID uint
	Amount    uint
	Price     float64
}

func (po *ProductOrder) FromString(in string) {
	split := strings.SplitN(in, ",", 3)
	if len(split) != 3 {
		log.Panicf("string '%s' is not a ProductOrder", in)
	}
	pid, err := strconv.ParseUint(split[0], 10, 64)
	if err != nil {
		log.Panicf("failed to parse string '%s' into a uint: %s", split[0], err)
	}
	amount, err := strconv.ParseUint(split[1], 10, 64)
	if err != nil {
		log.Panicf("failed to parse string '%s' into a uint: %s", split[1], err)
	}
	price, err := strconv.ParseFloat(split[2], 64)
	if err != nil {
		log.Panicf("failed to parse string '%s' into a float: %s", split[0], err)
	}
	po.ProductID = uint(pid)
	po.Amount = uint(amount)
	po.Price = price
}

func (po ProductOrder) String() string {
	return fmt.Sprintf("%d,%d,%f", po.ProductID, po.Amount, po.Price)
}

type OrderController interface {
	Register(Router)
}

type OrderUsecase interface {
	//
}

type OrderRepository interface {
	//
}
