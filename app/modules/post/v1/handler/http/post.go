package http

import (
	"errors"
	"strconv"

	"github.com/lichcse/clawathon-2026-3anhem-code/app/middleware"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/post/v1/service"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/post/v1/validation"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/response"
	schema "github.com/lichcse/clawathon-2026-3anhem-code/app/schema/post/v1"

	"github.com/gin-gonic/gin"
)

// PostHandler struct
type PostHandler struct {
	postService    service.PostService
	response       response.IMResponse
	postValidation validation.PostValidation
}

// NewPostHandler func new post handler
func NewPostHandler(
	postService service.PostService,
	response response.IMResponse,
	postValidation validation.PostValidation,
) *PostHandler {
	return &PostHandler{
		postService:    postService,
		response:       response,
		postValidation: postValidation,
	}
}

// Add func godoc
// @Summary Add a new post
// @Description Author: LichTV
// @Tags post
// @Accept json
// @Produce json
// @Param lang query string false "string" enums(en, vi)
// @Param PostAddRequest body schema.PostAddRequest true "Add a new post body"
// @Success 200 {object} schema.PostDetailResponse "success"
// @Router /post/v1/post [post]
func (h *PostHandler) Add(ctx *gin.Context) {
	// Get user ID from middleware
	userID, exists := ctx.Get(middleware.UserIDKey)
	if !exists {
		h.response.Out(ctx, errors.New("unauthorized"), nil)
		return
	}

	postAddRequest := &schema.PostAddRequest{}
	err := ctx.BindJSON(postAddRequest)
	if err != nil {
		h.response.Out(ctx, errors.New("not_allow"), nil)
		return
	}

	err = h.postValidation.Add(postAddRequest)
	if err != nil {
		h.response.Out(ctx, err, nil)
		return
	}

	postDetailResponse, err := h.postService.Add(userID.(uint64), postAddRequest)
	h.response.Out(ctx, err, postDetailResponse)
}

// Detail func godoc
// @Summary Detail info of post
// @Description Author: LichTV
// @Tags post
// @Accept json
// @Produce json
// @Param lang query string false "string" enums(en, vi)
// @Param id path int true "number"
// @Success 200 {object} schema.PostDetailResponse "success"
// @Router /post/v1/post/{id} [get]
func (h *PostHandler) Detail(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		h.response.Out(ctx, errors.New("not_allow"), nil)
		return
	}

	postDetailResponse, err := h.postService.Detail(uint64(id))
	h.response.Out(ctx, err, postDetailResponse)
}

// Update func godoc
// @Summary Update a post
// @Description Author: LichTV
// @Tags post
// @Accept json
// @Produce json
// @Param lang query string false "string" enums(en, vi)
// @Param id path int true "number"
// @Param PostUpdateRequest body schema.PostUpdateRequest true "Update post body"
// @Success 200 {object} response.Out "success"
// @Router /post/v1/post/{id} [put]
func (h *PostHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.response.Out(ctx, errors.New("not_allow"), nil)
		return
	}

	postUpdateRequest := &schema.PostUpdateRequest{}
	err := ctx.BindJSON(postUpdateRequest)
	if err != nil {
		h.response.Out(ctx, errors.New("not_allow"), nil)
		return
	}

	err = h.postValidation.Update(postUpdateRequest)
	if err != nil {
		h.response.Out(ctx, err, nil)
		return
	}

	err = h.postService.Update(id, postUpdateRequest)
	h.response.Out(ctx, err, nil)
}

// Delete func godoc
// @Summary Delete a post
// @Description Author: LichTV
// @Tags post
// @Accept json
// @Produce json
// @Param lang query string false "string" enums(en, vi)
// @Param id path int true "number"
// @Success 200 {object} response.Out "success"
// @Router /post/v1/post/{id} [delete]
func (h *PostHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.response.Out(ctx, errors.New("not_allow"), nil)
		return
	}

	err := h.postService.Delete(id)
	h.response.Out(ctx, err, nil)
}

// List func godoc
// @Summary List posts with pagination
// @Description Author: LichTV
// @Tags post
// @Accept json
// @Produce json
// @Param lang query string false "string" enums(en, vi)
// @Param page query int false "page number"
// @Param limit query int false "limit per page"
// @Success 200 {object} schema.PostListResponse "success"
// @Router /post/v1/post [get]
func (h *PostHandler) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Get user ID from query param (optional - for filtering by user)
	userIDStr := ctx.Query("user_id")
	var userID uint64
	if userIDStr != "" {
		userID, _ = strconv.ParseUint(userIDStr, 10, 64)
	}

	postListResponse, err := h.postService.List(userID, page, limit)
	h.response.Out(ctx, err, postListResponse)
}
