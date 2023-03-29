package domain

type Context interface {
	GetHeader(header string) string
	JSON(code int, i any)
	Status(code int)
	String(code int, format string, value ...any)
}
