package domain

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	. "github.com/minebarteksa/clean-show/logger"
)

type Order struct {
	DBModel         `json:"-"`
	Status          OrderStatus           `db:"status" json:"status"`
	OrderBy         uint                  `db:"order_by" json:"order_by" patch:"-"`
	ShippingAddress string                `db:"shipping_address" json:"shipping_address"`
	InvoiceAddress  string                `db:"invoice_address" json:"invoice_address"`
	Products        DBArray[ProductOrder] `db:"products" json:"products"`
	ShippingPrice   float64               `db:"shipping_price" json:"shipping_price"`
	Total           float64               `db:"total" json:"total" patch:"-"`
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
	ProductID uint    `json:"product_id"`
	Amount    uint    `json:"amount"`
	Price     float64 `json:"price"`
}

func (po *ProductOrder) FromString(in string) {
	if in == "" {
		return
	}
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

func (po *ProductOrder) FromInterfaceMap(in map[string]any) {
	po.ProductID = uint(jsonToInt(in["product_id"]))
	po.Amount = uint(jsonToInt(in["product_id"]))
	po.Price = jsonToFloat(in["price"])
}

func jsonToInt(in any) int {
	i, ok := in.(int64)
	if ok {
		return int(i)
	}
	f, ok := in.(float64)
	if ok {
		return int(f)
	}
	Log.Panicf("can not cast %T to int", in)
	return 0
}

func jsonToFloat(in any) float64 {
	f, ok := in.(float64)
	if ok {
		return f
	}
	i, ok := in.(int64)
	if ok {
		return float64(i)
	}
	Log.Panicf("can not cast %T to float64", in)
	return 0
}

func (po ProductOrder) String() string {
	return fmt.Sprintf("%d,%d,%f", po.ProductID, po.Amount, po.Price)
}

type OrderCreate struct {
	ShippingAddress string         `json:"shipping_address"`
	InvoiceAddress  string         `json:"invoice_address"`
	Products        []ProductOrder `json:"products"`
	ShippingPrice   float64        `json:"shipping_price"`
}

func (oc *OrderCreate) ToOrder(orderBy uint) *Order {
	order := &Order{
		Status:          OrderStatusCreated,
		OrderBy:         orderBy,
		ShippingAddress: oc.ShippingAddress,
		InvoiceAddress:  oc.InvoiceAddress,
		Products:        oc.Products,
		ShippingPrice:   oc.ShippingPrice,
	}
	order.UpdateTotal()
	return order
}

func (o *Order) UpdateTotal() {
	o.Total = float64(0)
	for _, product := range o.Products {
		o.Total += product.Price * float64(product.Amount)
	}
	o.Total += o.ShippingPrice
}

//go:generate mockery --name OrderController
type OrderController interface {
	Register(router Router)
	Get(context Context, session UserSession)
	Post(context Context, session UserSession)
	GetByID(context Context, session UserSession)
	Patch(context Context, session UserSession)
	PostCancel(context Context, session UserSession)
	Delete(context Context, session UserSession)
}

//go:generate mockery --name OrderUsecase
type OrderUsecase interface {
	TotalCount() (uint, error)
	Fetch(limit, page int) ([]Order, error)
	FetchByAccount(accountId uint, limit, page int) ([]Order, error)
	Create(accountId uint, create *OrderCreate) (*Order, error)
	FetchByID(session UserSession, id uint) (*Order, error)
	Modify(accountId, orderId uint, data map[string]any) (*Order, error)
	Cancel(session UserSession, orderId uint) error
	CancelByAccount(executorId, accountId uint) error
	Remove(accountId, orderId uint) error
}

//go:generate mockery --name OrderRepository
type OrderRepository interface {
	Count() (uint, error)
	Select(limit, page int) ([]Order, error)
	SelectAccount(accountId uint, limit, page int) ([]Order, error)
	SelectID(id uint) (*Order, error)
	SelectOrderBy(orderId uint) (uint, error)
	Insert(order *Order) error
	Update(order *Order) error
	UpdateStatus(orderId uint, status OrderStatus) error
	BatchUpdateStatus(orders []uint, status OrderStatus) error
	Delete(id uint) error
}
