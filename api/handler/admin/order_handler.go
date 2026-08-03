package admin

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/handler"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	handler.BaseHandler
}

func NewOrderHandler(app *core.AppServer, db *gorm.DB) *OrderHandler {
	return &OrderHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}}
}

// RegisterRoutes 注册路由
func (h *OrderHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/order/")

	// 需要管理员登录的接口
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.POST("list", h.List)
		rg.GET("remove", h.Remove)
		rg.GET("clear", h.Clear)
		rg.GET("recharge", h.Recharge)
	}
}

func (h *OrderHandler) List(c *gin.Context) {
	var data struct {
		OrderNo  string   `json:"order_no"`
		Status   int      `json:"status"`
		PayTime  []string `json:"pay_time"`
		Page     int      `json:"page"`
		PageSize int      `json:"page_size"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	session := h.DB.Session(&gorm.Session{})
	if data.OrderNo != "" {
		session = session.Where("order_no", data.OrderNo)
	}
	if len(data.PayTime) == 2 {
		start := utils.Str2stamp(data.PayTime[0] + " 00:00:00")
		end := utils.Str2stamp(data.PayTime[1] + " 23:59:59")
		session = session.Where("pay_time >= ? AND pay_time <= ?", start, end)
	}
	if data.Status >= 0 {
		session = session.Where("status", data.Status)
	}
	var total int64
	session.Model(&model.Order{}).Count(&total)
	var items []model.Order
	var list = make([]vo.Order, 0)
	offset := (data.Page - 1) * data.PageSize
	res := session.Order("id DESC").Offset(offset).Limit(data.PageSize).Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var order vo.Order
			err := utils.CopyObject(item, &order)
			if err == nil {
				order.Id = uint(item.Id)
				order.CreatedAt = item.CreatedAt.Unix()
				order.UpdatedAt = item.UpdatedAt.Unix()
				channel, ok := types.PayChannel[item.Channel]
				if !ok {
					channel = item.PayWay
				}
				payName, ok := types.PayWays[item.PayWay]
				if !ok {
					payName = item.PayWay
				}
				order.ChannelName = channel
				order.PayName = payName
				list = append(list, order)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, vo.NewPage(total, data.Page, data.PageSize, list))
}

func (h *OrderHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	if id > 0 {
		var item model.Order
		res := h.DB.First(&item, id)
		if res.Error != nil {
			resp.ERROR(c, "记录不存在！")
			return
		}

		if item.Status == types.OrderPaidSuccess {
			resp.ERROR(c, "已支付订单不允许删除！")
			return
		}

		err := h.DB.Where("id = ?", id).Delete(&model.Order{}).Error
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}
	resp.SUCCESS(c)
}

func (h *OrderHandler) Clear(c *gin.Context) {
	var orders []model.Order
	err := h.DB.Where("status <> ?", types.OrderPaidSuccess).Where("pay_time", 0).Find(&orders).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	deleteIds := make([]uint, 0)
	for _, order := range orders {
		// 只删除 15 分钟内的未支付订单
		if time.Now().After(order.CreatedAt.Add(time.Minute * 15)) {
			deleteIds = append(deleteIds, uint(order.Id))
		}
	}
	err = h.DB.Where("id IN ?", deleteIds).Delete(&model.Order{}).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}

// 补单
func (h *OrderHandler) Recharge(c *gin.Context) {
	id := h.GetInt(c, "id", 0)
	if id == 0 {
		resp.ERROR(c, "参数错误！")
		return
	}
	var order model.Order
	res := h.DB.First(&order, id)
	if res.Error != nil {
		resp.ERROR(c, "订单不存在！")
		return
	}
	if order.Status == types.OrderPaidSuccess {
		resp.ERROR(c, "订单已支付，不能补单！")
		return
	}
	order.Checked = false
	order.CreatedAt = time.Now()
	err := h.DB.Model(&order).Select("checked", "created_at").Updates(&order).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	resp.SUCCESS(c)
}
