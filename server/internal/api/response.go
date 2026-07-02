package api

type BaseResponse[T any] struct {
	Data T `json:"data"`
}
