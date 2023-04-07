package domain

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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
	total := float64(0)
	for _, product := range oc.Products {
		total += product.Price
	}
	total += oc.ShippingPrice
	return &Order{
		Status:          OrderStatusCreated,
		OrderBy:         orderBy,
		ShippingAddress: oc.ShippingAddress,
		InvoiceAddress:  oc.InvoiceAddress,
		Products:        oc.Products,
		ShippingPrice:   oc.ShippingPrice,
		Total:           total,
	}
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
	FetchByAccount(accountId uint) ([]Order, error)
	Create(accountId uint, create *OrderCreate) (*Order, error)
	FetchByID(session UserSession, id uint) (*Order, error)
	Modify(accountId, orderId uint, data map[string]any) error
	Cancel(session UserSession, orderId uint) error
	CancelByAccount(executorId, accountId uint) error
	Remove(accountId, orderId uint) error
}

//go:generate mockery --name OrderRepository
type OrderRepository interface {
	Count() (uint, error)
	Select(limit, page int) ([]Order, error)
	SelectAccount(accountId uint) ([]Order, error)
	SelectID(id uint) (*Order, error)
	SelectOrderBy(orderId uint) (uint, error)
	Insert(order *Order) error
	Update(order *Order) error
	UpdateStatus(orderId uint, status OrderStatus) error
	Delete(id uint) error
}
