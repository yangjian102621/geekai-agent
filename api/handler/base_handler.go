package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"errors"
	"fmt"
	"geekai/core"
	"geekai/core/types"
	"geekai/log"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"strings"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

var logger = log.GetLogger()

type BaseHandler struct {
	App     *core.AppServer
	DB      *gorm.DB
	Creator model.Creator
}

func (h *BaseHandler) GetTrim(c *gin.Context, key string) string {
	return strings.TrimSpace(c.Query(key))
}

func (h *BaseHandler) PostInt(c *gin.Context, key string, defaultValue int) int {
	return utils.IntValue(c.PostForm(key), defaultValue)
}

func (h *BaseHandler) GetInt(c *gin.Context, key string, defaultValue int) int {
	return utils.IntValue(c.Query(key), defaultValue)
}

func (h *BaseHandler) GetFloat(c *gin.Context, key string) float64 {
	return utils.FloatValue(c.Query(key))
}
func (h *BaseHandler) PostFloat(c *gin.Context, key string) float64 {
	return utils.FloatValue(c.PostForm(key))
}

func (h *BaseHandler) GetBool(c *gin.Context, key string) bool {
	return utils.BoolValue(c.Query(key))
}
func (h *BaseHandler) PostBool(c *gin.Context, key string) bool {
	return utils.BoolValue(c.PostForm(key))
}
func (h *BaseHandler) GetUserKey(c *gin.Context) string {
	userId, ok := c.Get(types.LoginUserID)
	if !ok {
		return ""
	}
	return fmt.Sprintf("users/%v", userId)
}

func (h *BaseHandler) GetLoginUserId(c *gin.Context) uint {
	userId, ok := c.Get(types.LoginUserID)
	if !ok {
		return 0
	}
	return uint(utils.IntValue(utils.InterfaceToString(userId), 0))
}

func (h *BaseHandler) GetLoginAdminId(c *gin.Context) uint {
	adminId, ok := c.Get(types.LoginAdminID)
	if !ok {
		return 0
	}
	return uint(utils.IntValue(utils.InterfaceToString(adminId), 0))
}

func (h *BaseHandler) IsLogin(c *gin.Context) bool {
	return h.GetLoginUserId(c) > 0
}

func (h *BaseHandler) GetLoginUser(c *gin.Context) (model.User, error) {
	value, exists := c.Get(types.LoginUserCache)
	if exists {
		return value.(model.User), nil
	}

	userId, ok := c.Get(types.LoginUserID)
	if !ok {
		return model.User{}, errors.New("user not login")
	}

	var user model.User
	res := h.DB.Where("id", userId).First(&user)
	// 更新缓存
	if res.Error == nil {
		c.Set(types.LoginUserCache, user)
	}
	return user, res.Error
}

// 获取当前创作者
func (h *BaseHandler) GetCurrentCreator(c *gin.Context) *model.Creator {
	if h.Creator.Id > 0 {
		return &h.Creator
	}

	creatorId := h.GetInt(c, "creator_id", 0)
	var creator model.Creator
	if creatorId == 0 {
		h.DB.Where("user_id = ?", h.GetLoginUserId(c)).First(&creator)
	} else {
		h.DB.Where("id = ?", creatorId).First(&creator)
	}

	if creator.Id == 0 || creator.Check != vo.CheckStatusPass {
		resp.ERROR(c, "您还不是创作者或创作者未审核通过")
		return nil
	}

	h.Creator = creator

	return &creator
}

func pushMessage(c *gin.Context, msgType string, content any) {
	c.SSEvent(msgType, content)
	c.Writer.Flush()
}
