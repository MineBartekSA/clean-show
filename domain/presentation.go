package domain

type DataList[T any] struct {
	Hits  uint `json:"hits"`
	Pages uint `json:"pages"`
	Data  []T  `json:"data"`
}
