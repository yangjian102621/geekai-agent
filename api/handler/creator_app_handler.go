package handler

import (
	"slices"

	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreatorAppHandler struct {
	BaseHandler
	appService *service.AppService
}

func NewCreatorAppHandler(app *core.AppServer, db *gorm.DB, appService *service.AppService) *CreatorAppHandler {
	return &CreatorAppHandler{
		BaseHandler: BaseHandler{App: app, DB: db},
		appService:  appService,
	}
}

func (h *CreatorAppHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/creator/")

	// 公开接口（无需登录）
	rg.GET("/:username/apps", h.GetCreatorApps)

	// 需要创作者登录的接口
	rg.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		rg.GET("/apps", h.List)
		rg.POST("/apps/save", h.SaveApp)
		rg.GET("/apps/remove", h.RemoveApp)
		rg.POST("/apps/enable", h.Enable)
	}
}

// GetApps 获取创作者的应用列表
func (h *CreatorAppHandler) List(c *gin.Context) {
	creator := h.GetCurrentCreator(c)
	if creator == nil {
		return
	}

	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 10)
	name := c.Query("name")
	cid := c.Query("cid")

	var apps []model.App
	var total int64

	query := h.DB.Model(&model.App{}).Where("creator_id = ?", creator.Id)
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if cid != "" {
		query = query.Where("cid = ?", cid)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&apps)

	// 获取分类
	catIds := []uint{}
	for _, app := range apps {
		catIds = append(catIds, app.Cid)
	}
	catIds = slices.Compact(catIds)

	var categories []model.AppCategory
	h.DB.Where("id IN (?)", catIds).Find(&categories)
	catMap := make(map[uint]model.AppCategory)
	for _, category := range categories {
		catMap[category.Id] = category
	}

	var items []vo.App
	for _, app := range apps {
		var item vo.App
		err := utils.CopyObject(app, &item)
		if err != nil {
			continue
		}
		item.Params = []vo.WorkflowParam{}
		if app.Configs != "" {
			if err = utils.JsonDecode(app.Configs, &item.Configs); err != nil {
				logger.Error(err)
			}
		}
		if app.Params != "" {
			if err = utils.JsonDecode(app.Params, &item.Params); err != nil {
				logger.Error(err)
				item.Params = []vo.WorkflowParam{}
			}
		}
		item.CreatedAt = app.CreatedAt.Unix()
		item.UpdatedAt = app.UpdatedAt.Unix()
		item.Cname = catMap[app.Cid].Name
		items = append(items, item)
	}
	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, items))
}

// SaveApp 创建/更新应用
func (h *CreatorAppHandler) SaveApp(c *gin.Context) {
	creator := h.GetCurrentCreator(c)
	if creator == nil {
		return
	}

	var data struct {
		Id      uint               `json:"id"`
		Name    string             `json:"name"`
		Type    string             `json:"type"`
		Enabled bool               `json:"enabled"`
		Score   int                `json:"score"`
		Icon    string             `json:"icon"`
		Summary string             `json:"summary"`
		Configs vo.AppConfig       `json:"configs"`
		Params  []vo.WorkflowParam `json:"params"`
		Cid     uint               `json:"cid"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, "参数错误："+err.Error())
		return
	}

	app := model.App{
		CreatorId: h.Creator.Id,
		Name:      data.Name,
		Type:      types.AppType(data.Type),
		Enabled:   data.Enabled,
		Score:     data.Score,
		Icon:      data.Icon,
		Summary:   data.Summary,
		Configs:   utils.JsonEncode(data.Configs),
		Params:    utils.JsonEncode(data.Params),
		Cid:       data.Cid,
		Check:     int8(vo.CheckStatusPending),
	}

	if data.Id > 0 {
		err := h.DB.Model(&model.App{}).
			Where("id = ?", data.Id).
			Select("name", "type", "enabled", "score", "icon", "summary", "configs", "params", "cid").
			Updates(app).Error
		if err != nil {
			resp.ERROR(c, "更新失败："+err.Error())
			return
		}
	} else {
		if err := h.DB.Create(&app).Error; err != nil {
			resp.ERROR(c, "创建失败："+err.Error())
			return
		}
	}

	resp.SUCCESS(c)
}

func (h *CreatorAppHandler) Enable(c *gin.Context) {
	var data struct {
		Id      uint `json:"id"`
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Model(&model.App{}).Where("id", data.Id).UpdateColumn("enabled", data.Enabled)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}
	resp.SUCCESS(c)
}

// RemoveApp 删除应用
func (h *CreatorAppHandler) RemoveApp(c *gin.Context) {
	creator := h.GetCurrentCreator(c)
	if creator == nil {
		return
	}

	id := h.GetInt(c, "id", 0)
	err := h.appService.RemoveApp(id, creator.Id)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{
		"message": "删除成功",
	})
}

// GetCreatorApps 获取创作者的公开应用列表（C端）
func (h *CreatorAppHandler) GetCreatorApps(c *gin.Context) {
	username := c.Param("username")
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 12)
	categoryId := h.GetInt(c, "category_id", 0)

	// 获取创作者id
	var creator model.Creator
	h.DB.Where("username = ? AND enabled = 1 AND `check` = 1", username).First(&creator)
	if creator.Id == 0 {
		resp.ERROR(c, "创作者不存在")
		return
	}

	var apps []model.App
	var total int64

	query := h.DB.Session(&gorm.Session{}).Model(&model.App{}).Where("creator_id = ? AND enabled = 1 AND `check` = 1", creator.Id)
	if categoryId > 0 {
		query = query.Where("cid = ?", categoryId)
	}
	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&apps)

	var items []vo.App
	for _, app := range apps {
		var item vo.App
		err := utils.CopyObject(app, &item)
		if err != nil {
			continue
		}
		item.CreatedAt = app.CreatedAt.Unix()
		item.UpdatedAt = app.UpdatedAt.Unix()
		items = append(items, item)
	}

	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, items))
}
