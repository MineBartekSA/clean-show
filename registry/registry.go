package registry

import (
	"github.com/minebarteksa/clean-show/domain"
	"github.com/minebarteksa/clean-show/infrastructure/router"
	"github.com/minebarteksa/clean-show/usecase/account"
	"github.com/minebarteksa/clean-show/usecase/audit"
	"github.com/minebarteksa/clean-show/usecase/product"
	"github.com/minebarteksa/clean-show/usecase/session"
	"go.uber.org/fx"
)

type Registry interface {
	Start()
}

type registry struct {
	db domain.DB

	app *fx.App
}

func NewRegistry(db domain.DB) Registry {
	return &registry{db: db}
}

func (r *registry) Start() {
	r.app = fx.New(
		fx.Supply(
			fx.Annotate(r.db, fx.As(new(domain.DB))),
		),
		fx.Provide(
			fx.Annotate(audit.NewAuditRepository, fx.As(new(domain.AuditRepository))),
			fx.Annotate(audit.NewAuditUsecase, fx.As(new(domain.AuditUsecase))),
			fx.Annotate(session.NewSessionRepository, fx.As(new(domain.SessionRepository))),
			fx.Annotate(session.NewSessionUsecase, fx.As(new(domain.SessionUsecase))),
			fx.Annotate(account.NewAccountRepository, fx.As(new(domain.AccountRepository))),
			fx.Annotate(account.NewAccountUsecase, fx.As(new(domain.AccountUsecase))),
			fx.Annotate(router.NewRouter, fx.As(new(domain.Router))),
			fx.Annotate(product.NewProductRepository, fx.As(new(domain.ProductRepository))),
			fx.Annotate(product.NewProductUsecase, fx.As(new(domain.ProductUsecase))),
			fx.Annotate(product.NewProductController, fx.As(new(domain.Controller)), fx.ResultTags(`group:"controllers"`)),
			fx.Annotate(router.NewAppController, fx.As(new(domain.Controller)), fx.ParamTags(`group:"controllers"`)),
		),
		fx.Invoke(func(r domain.Router, c domain.Controller) {
			c.Register(r)
			r.Run()
		}),
	)

	r.app.Run()
}
