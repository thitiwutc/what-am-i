package game

type BaseResponse[T any] struct {
	Data T `json:"data"`
}
