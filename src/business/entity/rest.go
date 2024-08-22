package entity

const (
	XRequestId = "x-request-id"
	RequestId  = "RequestId"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	IsError   bool   `json:"is_error"`
}
