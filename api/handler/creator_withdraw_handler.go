package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyr	ight 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

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

type WithdrawHandler struct {
	BaseHandler
}

func NewWithdrawHandler(app *core.AppServer, db *gorm.DB) *WithdrawHandler {
	return &WithdrawHandler{BaseHandler: BaseHandler{App: app, DB: db}}
}

func (h *WithdrawHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/creator/withdraws/")

	// 需要创作者登录的接口
	rg.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		rg.POST("/list", h.List)
	}
}

// List 获取提现申请列表
func (h *WithdrawHandler) List(c *gin.Context) {
	creator := h.GetCurrentCreator(c)
	if creator == nil {
		return
	}

	var data struct {
		Page      int    `json:"page" binding:"required"`
		PageSize  int    `json:"page_size" binding:"required"`
		Status    string `json:"status"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, "参数错误："+err.Error())
		return
	}

	var withdraws []model.CreatorWithdraw
	var total int64

	// 构建查询
	query := h.DB.Session(&gorm.Session{}).Where("creator_id = ?", creator.Id)
	// 筛选条件
	if data.Status != "" {
		query = query.Where("status = ?", data.Status)
	}

	if data.StartDate != "" {
		query = query.Where("created_at >= ?", data.StartDate)
	}
	if data.EndDate != "" {
		query = query.Where("created_at <= ?", data.EndDate)
	}

	// 统计总数
	query.Model(&model.CreatorWithdraw{}).Count(&total)

	// 分页查询
	query.Offset((data.Page - 1) * data.PageSize).Limit(data.PageSize).Order("id DESC").Find(&withdraws)

	var items []vo.CreatorWithdraw
	for _, withdraw := range withdraws {
		var item vo.CreatorWithdraw
		err := utils.CopyObject(withdraw, &item)
		if err != nil {
			continue
		}
		item.CreatedAt = withdraw.CreatedAt.Unix()
		item.UpdatedAt = withdraw.UpdatedAt.Unix()
		items = append(items, item)
	}
	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, items))
}
