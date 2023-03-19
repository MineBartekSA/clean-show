package domain

type Context interface {
	JSON(code int, i any)
	Status(code int)
}
