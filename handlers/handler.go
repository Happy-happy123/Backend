package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "StudentService/models"
    "StudentService/services"    
)
type StudentHandler struct {
	studentService *services.StudentService
}

func NewStudentHandler(studentService *services.StudentService) *StudentHandler {
	return &StudentHandler{studentService: studentService}
}
// 获取所有学生
func (h *StudentHandler) ListStudents(c *gin.Context) {
    student, err := h.studentService.ListStudents()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, student)
}

// 创建学生
func (h *StudentHandler) CreateStudent(c *gin.Context) {
    var newStudent models.Student
    var err error
    if err := c.ShouldBindJSON(&newStudent); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if newStudent.Name == "" || newStudent.Age <= 0 || newStudent.Gender == "" {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
     _, err = h.studentService.CreateStudent(&newStudent)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, newStudent)
}

// 获取单个学生
func (h *StudentHandler) GetStudent(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的学生ID"})
        return
    }
    student, err := h.studentService.GetStudent(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if student == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "学生不存在"})
        return
    }
    c.JSON(http.StatusOK, student)
}

// 更新学生
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的学生ID"})
        return
    }
    var updateInfo models.StudentUpdate
    if err := c.ShouldBindJSON(&updateInfo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if updateInfo.Name == nil && updateInfo.Age == nil && updateInfo.Gender == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "没有提供更新字段"})
        return
    }
    if err := h.studentService.UpdateStudent(id, &updateInfo); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}

// 删除学生
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的学生ID"})
        return
    }
    if err := h.studentService.DeleteStudent(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}