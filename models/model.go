package models

type Student struct {
	Id     int    `json:"id"`
	Name   string `json:"name" binding:"required"`//只要请求 JSON 里 没有这个字段或值为零值，Gin 立刻拦下来返回 400，省掉一堆校验
	Age    int    `json:"age"`
	Gender string `json:"gender"`
	// 加了 json:"xxx" 后，HTTP 请求/响应可以直接 ShouldBind 绑定 JSON，或 c.JSON 输出 JSON
}
type StudentUpdate struct {
	Name   *string `json:"name"`
	Age    *int    `json:"age"`
	Gender *string `json:"gender"`
}