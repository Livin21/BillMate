package util

type MessageResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type DataResponse struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}