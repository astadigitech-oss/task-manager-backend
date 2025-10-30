package utils

type APIResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Optionally, helper function
func NewResponse(success bool, code int, msg string, data interface{}) APIResponse {
	return APIResponse{
		Success: success,
		Code:    code,
		Message: msg,
		Data:    data,
	}
}
