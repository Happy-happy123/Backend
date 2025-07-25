package student

type Student struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

// 请求和响应DTO
type CreateRequest struct {
	Name   string `json:"name" binding:"required"`
	Age    int    `json:"age" binding:"required,min=1"`
	Gender string `json:"gender" binding:"required,oneof=男 女"`
}

type UpdateRequest struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}