package handler

import (
	"geekai/core"
	"geekai/core/middleware"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppCategoryHandler struct {
	BaseHandler
}

func NewAppCategoryHandler(app *core.AppServer, db *gorm.DB) *AppCategoryHandler {
	return &AppCategoryHandler{
		BaseHandler: BaseHandler{App: app, DB: db},
	}
}

func (h *AppCategoryHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/creator/app-categories/")
	rg.GET("list", h.List)
	// 需要创作者登录的接口
	rg.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		rg.POST("/create", h.Create)
		rg.POST("/update", h.Update)
		rg.POST("/delete", h.Delete)
	}
}

// List 获取分类列表
func (h *AppCategoryHandler) List(c *gin.Context) {
	creator := h.GetCurrentCreator(c)
	if creator == nil {
		return
	}

	var categories []model.AppCategory
	query := h.DB.Session(&gorm.Session{}).Where("creator_id = ?", creator.Id)
	// 参数过滤
	enabled := c.Query("enabled")
	if enabled != "" {
		query = query.Where("enabled = ?", enabled)
	}
	query.Order("id DESC").Find(&categories)
	// 转换为VO
	var items []vo.AppCategory
	for _, category := range categories {
		var item vo.AppCategory
		err := utils.CopyObject(category, &item)
		if err != nil {
			continue
		}
		item.CreatedAt = category.CreatedAt.Unix()
		item.UpdatedAt = category.UpdatedAt.Unix()
		items = append(items, item)
	}
	resp.SUCCESS(c, items)
}

// Create 新增分类
func (h *AppCategoryHandler) Create(c *gin.Context) {
	creator := h.GetCurrentCreator(c)
	if creator == nil {
		return
	}

	var req struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, "参数错误："+err.Error())
		return
	}
	if req.Name == "" {
		resp.ERROR(c, "分类名称不能为空")
		return
	}
	logger.Infof("创建分类: %d, %+v", creator.Id, req)
	category := model.AppCategory{
		CreatorId: creator.Id,
		Name:      req.Name,
		Enabled:   req.Enabled,
	}
	if err := h.DB.Create(&category).Error; err != nil {
		resp.ERROR(c, "创建失败："+err.Error())
		return
	}
	resp.SUCCESS(c)
}

// Update 编辑分类
func (h *AppCategoryHandler) Update(c *gin.Context) {
	creator := h.GetCurrentCreator(c)
	if creator == nil {
		return
	}

	var req struct {
		Id      uint   `json:"id"`
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, "参数错误："+err.Error())
		return
	}
	if req.Id == 0 {
		resp.ERROR(c, "ID不能为空")
		return
	}

	if !req.Enabled {
		// 检查分类下是否有应用
		var count int64
		h.DB.Model(&model.App{}).Where("cid = ? AND enabled = 1", req.Id).Count(&count)
		if count > 0 {
			resp.ERROR(c, "应用分类下有启用的应用，不能禁用")
			return
		}
	}

	if err := h.DB.Model(&model.AppCategory{}).Where("id = ?", req.Id).Updates(map[string]any{
		"name":    req.Name,
		"enabled": req.Enabled,
	}).Error; err != nil {
		resp.ERROR(c, "更新失败："+err.Error())
		return
	}
	resp.SUCCESS(c)
}

// Delete 删除分类
func (h *AppCategoryHandler) Delete(c *gin.Context) {
	var req struct {
		Id uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.ERROR(c, "参数错误："+err.Error())
		return
	}
	if req.Id == 0 {
		resp.ERROR(c, "ID不能为空")
		return
	}
	// 检查分类下是否有应用
	var count int64
	h.DB.Model(&model.App{}).Where("cid = ?", req.Id).Count(&count)
	if count > 0 {
		resp.ERROR(c, "分类下有应用，不能删除")
		return
	}

	if err := h.DB.Delete(&model.AppCategory{}, req.Id).Error; err != nil {
		resp.ERROR(c, "删除失败："+err.Error())
		return
	}
	resp.SUCCESS(c, gin.H{"message": "删除成功"})
}
