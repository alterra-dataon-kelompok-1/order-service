package response

type BaseResponse struct {
	Meta
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}

// func SuccessResponse(c echo.Context, code int, data interface{}) error {
// 	response := BaseResponse{}
//
// 	return c.JSON(code, response)
// }
