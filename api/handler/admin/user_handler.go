package admin

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
	"geekai/handler"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	handler.BaseHandler
	redis      *redis.Client
	baseConfig types.BaseConfig
}

func NewUserHandler(app *core.AppServer, db *gorm.DB, redisCli *redis.Client, sysConfig *types.SystemConfig) *UserHandler {
	return &UserHandler{BaseHandler: handler.BaseHandler{App: app, DB: db}, redis: redisCli, baseConfig: sysConfig.Base}
}

// RegisterRoutes 注册路由
func (h *UserHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/user/")

	// 需要管理员登录的接口
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.GET("list", h.List)
		rg.POST("save", h.Save)
		rg.POST("enable", h.Enable)
		rg.GET("remove", h.Remove)
		rg.GET("loginLog", h.LoginLog)
		rg.POST("resetPass", h.ResetPass)
		rg.GET("nickname", h.GenerateNickname)
	}
}

// List 用户列表
func (h *UserHandler) List(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	username := h.GetTrim(c, "username")
	userId := h.GetInt(c, "user_id", 0)

	offset := (page - 1) * pageSize
	var items []model.User
	var users = make([]vo.User, 0)
	var total int64

	session := h.DB.Session(&gorm.Session{})
	if username != "" {
		session = session.Where("username LIKE ?", "%"+username+"%")
	}
	if userId > 0 {
		session = session.Where("id = ?", userId)
	}

	session.Model(&model.User{}).Count(&total)
	res := session.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var user vo.User
			err := utils.CopyObject(item, &user)
			if err == nil {
				user.Id = item.Id
				user.CreatedAt = item.CreatedAt.Unix()
				user.UpdatedAt = item.UpdatedAt.Unix()
				users = append(users, user)
			} else {
				logger.Error(err)
			}
		}
	}
	pageVo := vo.NewPage(total, page, pageSize, users)
	resp.SUCCESS(c, pageVo)
}

func (h *UserHandler) Save(c *gin.Context) {
	var data struct {
		Id          uint   `json:"id"`
		Password    string `json:"password"`
		Username    string `json:"username"`
		Nickname    string `json:"nickname"`
		ExpiredTime string `json:"expired_time"`
		Enabled     bool   `json:"enabled"`
		Vip         bool   `json:"vip"`
		Scores      int    `json:"scores"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var user = model.User{}
	var res *gorm.DB
	if data.Id > 0 { // 更新
		res = h.DB.Where("id", data.Id).First(&user)
		if res.Error != nil {
			resp.ERROR(c, "user not found")
			return
		}
		var oldScores = user.Scores
		user.Username = data.Username
		user.Nickname = data.Nickname
		user.Enabled = data.Enabled
		user.Vip = data.Vip
		user.Scores = data.Scores
		user.ExpiredTime = utils.Str2stamp(data.ExpiredTime)

		res = h.DB.Select("username", "nickname", "enabled", "vip", "scores", "expired_time").Updates(&user)

		if res.Error != nil {
			logger.Error("error with update database：", res.Error)
			resp.ERROR(c, res.Error.Error())
			return
		}
		// 记录算力日志
		if oldScores != user.Scores {
			mark := types.ScorePlus
			amount := user.Scores - oldScores
			if oldScores > user.Scores {
				mark = types.ScoreSub
				amount = oldScores - user.Scores
			}
			h.DB.Create(&model.ScoreLog{
				UserId:    user.Id,
				Username:  user.Username,
				Type:      types.ScoreFineTune,
				Amount:    amount,
				Balance:   user.Scores,
				Mark:      mark,
				Subject:   "管理员",
				Remark:    fmt.Sprintf("后台管理员强制修改用户算力，修改前：%d,修改后:%d, 管理员ID：%d", oldScores, user.Scores, h.GetLoginUserId(c)),
				CreatedAt: time.Now(),
			})
		}
		// 如果禁用了用户，则将用户踢下线
		if !user.Enabled {
			key := fmt.Sprintf("users/%v", user.Id)
			if _, err := h.redis.Del(c, key).Result(); err != nil {
				logger.Error("error with delete session: ", err)
			}
		}
	} else {
		// 检查用户是否已经存在
		h.DB.Where("username", data.Username).First(&user)
		if user.Id > 0 {
			resp.ERROR(c, "用户名已存在")
			return
		}

		salt := utils.RandString(8)
		u := model.User{
			Username:    data.Username,
			Nickname:    data.Nickname,
			Password:    utils.GenPassword(data.Password, salt),
			Salt:        salt,
			Scores:      data.Scores,
			Enabled:     true,
			ExpiredTime: utils.Str2stamp(data.ExpiredTime),
			Avatar:      "/images/avatar/user.png",
		}

		res = h.DB.Create(&u)
	}

	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}

	resp.SUCCESS(c)
}

// Enable 启用/禁用
func (h *UserHandler) Enable(c *gin.Context) {
	var data struct {
		Id      uint `json:"id"`
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.DB.Model(&model.User{}).Where("id", data.Id).UpdateColumn("enabled", data.Enabled)
	if res.Error != nil {
		resp.ERROR(c, res.Error.Error())
		return
	}
	resp.SUCCESS(c)
}

// ResetPass 重置密码
func (h *UserHandler) ResetPass(c *gin.Context) {
	var data struct {
		Id       uint
		Password string
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var user model.User
	res := h.DB.First(&user, data.Id)
	if res.Error != nil {
		resp.ERROR(c, "No user found")
		return
	}

	password := utils.GenPassword(data.Password, user.Salt)
	user.Password = password
	res = h.DB.Updates(&user)
	if res.Error != nil {
		resp.ERROR(c)
	} else {
		resp.SUCCESS(c)
	}
}

func (h *UserHandler) Remove(c *gin.Context) {
	id := c.Query("id")
	ids := c.QueryArray("ids[]")
	if id != "" {
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	tx := h.DB.Begin()
	var err error
	for _, id = range ids {
		// 删除用户
		if err = tx.Where("id", id).Delete(&model.User{}).Error; err != nil {
			break
		}
		// 删除算力日志
		if err = tx.Where("user_id = ?", id).Delete(&model.ScoreLog{}).Error; err != nil {
			break
		}
	}
	if err != nil {
		resp.ERROR(c, err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
	resp.SUCCESS(c)
}

func (h *UserHandler) LoginLog(c *gin.Context) {
	page := h.GetInt(c, "page", 1)
	pageSize := h.GetInt(c, "page_size", 20)
	var total int64
	h.DB.Model(&model.UserLoginLog{}).Count(&total)
	offset := (page - 1) * pageSize
	var items []model.UserLoginLog
	res := h.DB.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
	if res.Error != nil {
		resp.ERROR(c, "获取数据失败")
		return
	}
	userIds := make([]uint, 0)
	for _, v := range items {
		userIds = append(userIds, v.UserId)
	}
	users := make([]model.User, 0)
	h.DB.Where("id IN (?)", userIds).Find(&users)
	userMap := make(map[uint]model.User)
	for _, v := range users {
		userMap[v.Id] = v
	}
	var logs []vo.UserLoginLog
	for _, v := range items {
		var log vo.UserLoginLog
		err := utils.CopyObject(v, &log)
		if err == nil {
			log.Id = v.Id
			log.CreatedAt = v.CreatedAt.Unix()
			if user, ok := userMap[v.UserId]; ok {
				log.Nickname = user.Nickname
			}
			logs = append(logs, log)
		}
	}

	resp.SUCCESS(c, vo.NewPage(total, page, pageSize, logs))
}

// 生成用户昵称
func (h *UserHandler) GenerateNickname(c *gin.Context) {
	nickname := utils.GenerateNickname()
	resp.SUCCESS(c, gin.H{
		"nickname": nickname,
	})
}
