package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yang/go-learning-backend/gin-backend/internal/errs"
	"github.com/yang/go-learning-backend/gin-backend/internal/response"
	"github.com/yang/go-learning-backend/gin-backend/internal/service"
	"github.com/yang/go-learning-backend/gin-backend/internal/transport/http/dto"
	"github.com/yang/go-learning-backend/gin-backend/internal/transport/http/middleware"
)

// TodoHandler 处理 Todo 接口。
type TodoHandler struct {
	todoService service.TodoService
}

// NewTodoHandler 创建 Todo 处理器。
func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

// Create 创建 Todo。
func (h *TodoHandler) Create(c *gin.Context) {
	// user_id 由 JWT 中间件注入。
	userID := c.GetString(middleware.UserIDKey())
	if userID == "" {
		response.Error(c, errs.WithMessage(errs.ErrUnauthorized, "missing user identity"))
		return
	}
	var req dto.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, responseBindError(err))
		return
	}

	created, err := h.todoService.Create(c.Request.Context(), userID, req.Title)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, created)
}

// List 返回当前登录用户的 Todo 列表。
func (h *TodoHandler) List(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey())
	if userID == "" {
		response.Error(c, errs.WithMessage(errs.ErrUnauthorized, "missing user identity"))
		return
	}
	items, err := h.todoService.List(c.Request.Context(), userID)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, http.StatusOK, items)
}

// Update 更新指定 Todo。
func (h *TodoHandler) Update(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey())
	if userID == "" {
		response.Error(c, errs.WithMessage(errs.ErrUnauthorized, "missing user identity"))
		return
	}
	id := c.Param("id")

	var req dto.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, responseBindError(err))
		return
	}

	updated, err := h.todoService.Update(c.Request.Context(), userID, id, service.UpdateTodoInput{
		Title: req.Title,
		Done:  req.Done,
	})
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, http.StatusOK, updated)
}

// Delete 删除指定 Todo。
func (h *TodoHandler) Delete(c *gin.Context) {
	userID := c.GetString(middleware.UserIDKey())
	if userID == "" {
		response.Error(c, errs.WithMessage(errs.ErrUnauthorized, "missing user identity"))
		return
	}
	id := c.Param("id")
	if err := h.todoService.Delete(c.Request.Context(), userID, id); err != nil {
		writeServiceError(c, err)
		return
	}
	response.NoContent(c)
}
