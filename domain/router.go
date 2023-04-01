package domain

type Handler func(context Context, session UserSession)

type Router interface {
	Run()
	API() RouteGroup
	Auth(token string) (*Session, *Account, error)
}

type RouteGroup interface {
	Group(relativeRath string) RouteGroup
	GET(relativePath string, handlers Handler, authorized AuthLevel)
	POST(relativePath string, handlers Handler, authorized AuthLevel)
	PATCH(relativePath string, handlers Handler, authorized AuthLevel)
	DELETE(relativePath string, handlers Handler, authorized AuthLevel)
}

type AuthLevel int

// Should map directly to AccountType
const (
	AuthLevelNone AuthLevel = iota
	AuthLevelUser
	AuthLevelStaff
)

type Controller interface {
	Register(router Router)
}
