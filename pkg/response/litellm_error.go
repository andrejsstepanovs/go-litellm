package response

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
	Param   string `json:"param"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}
