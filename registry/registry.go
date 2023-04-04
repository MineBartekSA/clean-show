package registry

import (
	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/infrastructure/router"
	"github.com/minebarteksa/clean-show/usecase/account"
	"github.com/minebarteksa/clean-show/usecase/audit"
	"github.com/minebarteksa/clean-show/usecase/order"
	"github.com/minebarteksa/clean-show/usecase/product"
	"github.com/minebarteksa/clean-show/usecase/session"
	"go.uber.org/fx"
)

type Registry interface {
	Start()
}

type registry struct {
	db     domain.DB
	hasher domain.Hasher
}

func NewRegistry(db domain.DB, hasher domain.Hasher) Registry {
	return &registry{db, hasher}
}

func (r *registry) Start() {
	app := fx.New(
		fx.Supply(
			fx.Annotate(r.db, fx.As(new(domain.DB))),
			fx.Annotate(r.hasher, fx.As(new(domain.Hasher))),
		),
		fx.Provide(
			fx.Annotate(audit.NewAuditRepository, fx.As(new(domain.AuditRepository))),
			fx.Annotate(audit.NewAuditUsecase, fx.As(new(domain.AuditUsecase))),
			fx.Annotate(session.NewSessionRepository, fx.As(new(domain.SessionRepository))),
			fx.Annotate(session.NewSessionUsecase, fx.As(new(domain.SessionUsecase))),
			fx.Annotate(product.NewProductRepository, fx.As(new(domain.ProductRepository))),
			fx.Annotate(product.NewProductUsecase, fx.As(new(domain.ProductUsecase))),
			fx.Annotate(product.NewProductController, fx.As(new(domain.Controller)), fx.ResultTags(`group:"controllers"`)),
			fx.Annotate(order.NewOrderRepository, fx.As(new(domain.OrderRepository))),
			fx.Annotate(order.NewOrderUsecase, fx.As(new(domain.OrderUsecase))),
			fx.Annotate(order.NewOrderController, fx.As(new(domain.Controller)), fx.ResultTags(`group:"controllers"`)),
			fx.Annotate(account.NewAccountRepository, fx.As(new(domain.AccountRepository))),
			fx.Annotate(account.NewAccountUsecase, fx.As(new(domain.AccountUsecase))),
			fx.Annotate(account.NewAccountController, fx.As(new(domain.Controller)), fx.ResultTags(`group:"controllers"`)),
			fx.Annotate(router.NewRouter, fx.As(new(domain.Router))),
			fx.Annotate(router.NewAppController, fx.As(new(domain.Controller)), fx.ParamTags(`group:"controllers"`)),
		),
		fx.Invoke(func(r domain.Router, c domain.Controller) {
			c.Register(r)
			r.Run()
		}),
	)

	app.Run()
}
