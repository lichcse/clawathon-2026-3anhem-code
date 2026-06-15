package http

import (
	"errors"
	"strconv"

	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/identity/v1/service"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/modules/identity/v1/validation"
	"github.com/lichcse/clawathon-2026-3anhem-code/app/response"
	schema "github.com/lichcse/clawathon-2026-3anhem-code/app/schema/identity/v1"

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
// @Tags identity
// @Accept json
// @Produce json
// @Param lang query string false "string" enums(en, vi)
// @Param X-User-ID header int false "User ID from auth middleware"
// @Param PostAddRequest body schema.PostAddRequest true "Add a new post body"
// @Success 200 {object} schema.PostDetailResponse "success"
// @Router /identity/v1/post [post]
func (p *PostHandler) Add(ctx *gin.Context) {
	// Get userID from header (simulated auth)
	userIDStr := ctx.GetHeader("X-User-ID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		p.response.Out(ctx, errors.New("not_allow"), nil)
		return
	}

	postAddRequest := &schema.PostAddRequest{}
	err = ctx.BindJSON(postAddRequest)
	if err != nil {
		p.response.Out(ctx, errors.New("not_allow"), nil)
		return
	}

	err = p.postValidation.Add(postAddRequest)
	if err != nil {
		p.response.Out(ctx, err, nil)
		return
	}

	postDetailResponse, err := p.postService.Add(userID, postAddRequest)
	if err != nil {
		p.response.Out(ctx, err, nil)
		return
	}
	p.response.Out(ctx, err, postDetailResponse)
}

// Detail func godoc
// @Summary Detail info of post
// @Description Author: LichTV
// @Tags identity
// @Accept json
// @Produce json
// @Param lang query string false "string" enums(en, vi)
// @Param id path int true "number"
// @Success 200 {object} schema.PostDetailResponse "success"
// @Router /identity/v1/post/{id} [get]
func (p *PostHandler) Detail(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		p.response.Out(ctx, errors.New("not_allow"), nil)
		return
	}

	postDetailResponse, err := p.postService.Detail(uint64(id))
	p.response.Out(ctx, err, postDetailResponse)
	return
}

// Update func godoc
// @Summary Update a post
// @Description Author: LichTV
// @Tags identity
// @Accept json
// @Produce json
// @Param lang query string false "string" enums(en, vi)
// @Param X-User-ID header int false "User ID from auth middleware"
// @Param id path int true "number"
// @Param PostUpdateRequest body schema.PostUpdateRequest true "Update post body"
// @Success 200 {object} schema.PostDetailResponse "success"
// @Router /identity/v1/post/{id} [put]
func (p *PostHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		p.response.Out(ctx, errors.New("not_allow"), nil)
		return
	}

	postUpdateRequest := &schema.PostUpdateRequest{}
	err := ctx.BindJSON(postUpdateRequest)
	if err != nil {
		p.response.Out(ctx, errors.New("not_allow"), nil)
		return
	}

	err = p.postValidation.Update(postUpdateRequest)
	if err != nil {
		p.response.Out(ctx, err, nil)
		return
	}

	err = p.postService.Update(id, postUpdateRequest)
	p.response.Out(ctx, err, nil)
}

// Delete func godoc
// @Summary Delete a post
// @Description Author: LichTV
// @Tags identity
// @Accept json
// @Produce json
// @Param lang query string false "string" enums(en, vi)
// @Param id path int true "number"
// @Success 200 {object} schema.PostDetailResponse "success"
// @Router /identity/v1/post/{id} [delete]
func (p *PostHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		p.response.Out(ctx, errors.New("not_allow"), nil)
		return
	}

	err := p.postService.Delete(id)
	p.response.Out(ctx, err, nil)
}

// ListByUserID func godoc
// @Summary List posts by user ID
// @Description Author: LichTV
// @Tags identity
// @Accept json
// @Produce json
// @Param lang query string false "string" enums(en, vi)
// @Param X-User-ID header int false "User ID from auth middleware"
// @Success 200 {object} []schema.PostDetailResponse "success"
// @Router /identity/v1/post/user [get]
func (p *PostHandler) ListByUserID(ctx *gin.Context) {
	// Get userID from header (simulated auth)
	userIDStr := ctx.GetHeader("X-User-ID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		p.response.Out(ctx, errors.New("not_allow"), nil)
		return
	}

	posts, err := p.postService.ListByUserID(userID)
	p.response.Out(ctx, err, posts)
}