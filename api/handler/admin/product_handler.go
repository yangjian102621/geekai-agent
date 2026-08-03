package admin

import (
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/handler"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	handler.BaseHandler
}

func NewProductHandler(app *core.AppServer, db *gorm.DB) *ProductHandler {
	return &ProductHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

// RegisterRoutes 注册路由
func (h *ProductHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/product/")

	// 需要管理员登录的接口
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.GET("list", h.List)
		rg.POST("save", h.Save)
		rg.POST("enable", h.Enable)
		rg.GET("remove", h.Remove)
		rg.POST("sort", h.Sort)
	}
}

func (h *ProductHandler) Save(c *gin.Context) {
	var data vo.Product
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	item := model.Product{
		Name:    data.Name,
		Price:   data.Price,
		Credit:  data.Credit,
		Enabled: data.Enabled,
	}
	item.Id = data.Id
	var err error
	if item.Id > 0 {
		err = h.DB.Model(&model.Product{}).Where("id", item.Id).Updates(map[string]any{
			"name":    item.Name,
			"price":   item.Price,
			"credit":  item.Credit,
			"enabled": item.Enabled,
		}).Error
	} else {
		err = h.DB.Create(&item).Error
	}

	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

// List 数据列表
func (h *ProductHandler) List(c *gin.Context) {
	var items []model.Product
	var list = make([]vo.Product, 0)
	res := h.DB.Order("sort_num ASC").Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var product vo.Product
			err := utils.CopyObject(item, &product)
			if err == nil {
				product.Id = item.Id
				product.CreatedAt = item.CreatedAt.Unix()
				product.UpdatedAt = item.UpdatedAt.Unix()
				list = append(list, product)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, list)
}

func (h *ProductHandler) Enable(c *gin.Context) {
	var data struct {
		Id      uint `json:"id"`
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	err := h.DB.Model(&model.Product{}).Where("id", data.Id).UpdateColumn("enabled", data.Enabled).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

func (h *ProductHandler) Sort(c *gin.Context) {
	var data struct {
		Ids   []uint `json:"ids"`
		Sorts []int  `json:"sorts"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	for index, id := range data.Ids {
		err := h.DB.Model(&model.Product{}).Where("id", id).Update("sort_num", data.Sorts[index]).Error
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}

	resp.SUCCESS(c)
}

func (h *ProductHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	if id > 0 {
		err := h.DB.Where("id", id).Delete(&model.Product{}).Error
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}
	resp.SUCCESS(c)
}
