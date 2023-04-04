package domain

type Context interface {
	Param(key string) string
	Query(key string) string
	GetHeader(header string) string
	UnmarshalBody(in any) error
	Cookie(name string) (string, error)

	JSON(code int, i any)
	Status(code int)
	String(code int, format string, value ...any)
	SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool)
}
