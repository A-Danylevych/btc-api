package endpoints

type GetRequest struct{}

type GetResponse struct {
	Rate float64 `json:"rate"`
	Err  string  `json:"err,omitempty"`
}

type ServiceStatusRequest struct{}

type ServiceStatusResponse struct {
	Code int    `json:"status"`
	Err  string `json:"err,omitempty"`
}
