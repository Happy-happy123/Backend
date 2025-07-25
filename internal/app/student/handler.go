package student

import (
	"strconv"
	
	"StudentService/pkg/api"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}//创建一个新的Handler结构体实例，并将传入的svc参数赋值给Handler的svc字段,返回实例的地址实际是指向Handler的一个指针
	}

func (h *Handler) CreateStudent(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "无效的请求数据")
		return
	}

	student, err := h.svc.CreateStudent(&req)
	if err != nil {
		api.ServerError(c, "创建学生失败")
		return
	}

	api.Created(c, student)
}

func (h *Handler) GetStudent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	student, err := h.svc.GetStudent(id)
	if err != nil {
		if err == ErrNotFound {
			api.NotFound(c, "学生不存在")
			return
		}
		api.ServerError(c, "获取学生失败")
		return
	}

	api.Success(c, student)
}

func (h *Handler) ListStudents(c *gin.Context) {
	students, err := h.svc.GetAllStudents()
	if err != nil {
		api.ServerError(c, "获取学生列表失败")
		return
	}
	api.Success(c, students)
}

func (h *Handler) UpdateStudent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.BadRequest(c, "无效的请求数据")
		return
	}

	student, err := h.svc.UpdateStudent(id, &req)
	if err != nil {
		if err == ErrNotFound {
			api.NotFound(c, "学生不存在")
			return
		}
		api.ServerError(c, "更新学生失败")
		return
	}

	api.Success(c, student)
}

func (h *Handler) DeleteStudent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.DeleteStudent(id); err != nil {
		if err == ErrNotFound {
			api.NotFound(c, "学生不存在")
			return
		}
		api.ServerError(c, "删除学生失败")
		return
	}

	api.Success(c, gin.H{"message": "学生已删除"})
}