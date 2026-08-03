package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"geekai/store/model"
	"geekai/utils/resp"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RedeemHandler struct {
	BaseHandler
	lock        sync.Mutex
	userService *service.UserService
}

func NewRedeemHandler(app *core.AppServer, db *gorm.DB, userService *service.UserService) *RedeemHandler {
	return &RedeemHandler{BaseHandler: BaseHandler{App: app, DB: db}, userService: userService}
}

// RegisterRoutes 注册路由
func (h *RedeemHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/redeem/")

	// 需要用户登录的接口
	rg.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		rg.POST("verify", h.Verify)
	}
}

func (h *RedeemHandler) Verify(c *gin.Context) {
	var data struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}
	userId := h.GetLoginUserId(c)

	h.lock.Lock()
	defer h.lock.Unlock()

	var item model.Redeem
	res := h.DB.Where("code", data.Code).First(&item)
	if res.Error != nil {
		resp.ERROR(c, "无效的兑换码！")
		return
	}

	if !item.Enabled {
		resp.ERROR(c, "当前兑换码已被禁用！")
		return
	}

	if item.RedeemedAt > 0 {
		resp.ERROR(c, "当前兑换码已使用，请勿重复使用！")
		return
	}

	tx := h.DB.Begin()
	err := h.userService.IncreaseScores(int(userId), int(item.Amount), model.ScoreLog{
		Type:    types.ScoreRedeem,
		Subject: "兑换码",
		Remark:  fmt.Sprintf("兑换码核销，额度：%d，兑换码：%s...", item.Amount, item.Code[:10]),
	})
	if err != nil {
		tx.Rollback()
		resp.ERROR(c, err.Error())
		return
	}

	// 更新核销状态
	item.RedeemedAt = time.Now().Unix()
	item.UserId = userId
	err = tx.Updates(&item).Error
	if err != nil {
		tx.Rollback()
		resp.ERROR(c, err.Error())
		return
	}

	tx.Commit()
	resp.SUCCESS(c, item.Amount)

}
