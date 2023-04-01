package domain

type Context interface {
	Param(key string) string
	Query(key string) string
	GetHeader(header string) string
	UnmarshalBody(in any) error

	JSON(code int, i any)
	Status(code int)
	String(code int, format string, value ...any)
}
