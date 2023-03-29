package domain

type Handler func(context Context, session UserSession)

type Router interface {
	Run()
	API() RouteGroup
	Auth(token string) (*Session, *Account, error)
}

type RouteGroup interface {
	Group(relativeRath string) RouteGroup
	GET(relativePath string, handlers Handler, authorized AccountType)
	POST(relativePath string, handlers Handler, authorized AccountType)
	PATCH(relativePath string, handlers Handler, authorized AccountType)
	DELETE(relativePath string, handlers Handler, authorized AccountType)
}

type Controller interface {
	Register(router Router)
}
