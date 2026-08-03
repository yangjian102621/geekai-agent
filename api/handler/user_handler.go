package handler

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"fmt"
	"geekai/core"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/service"
	"geekai/store/model"
	"geekai/store/vo"
	"geekai/utils"
	"geekai/utils/resp"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	BaseHandler
	redis          *redis.Client
	captcha        *service.CaptchaService
	userService    *service.UserService
	baseConfig     *types.BaseConfig
	wxLoginConfig  *types.WxLoginConfig
	wxLoginService *service.WxLoginService
	ipSearcher     *xdb.Searcher
}

func NewUserHandler(
	app *core.AppServer,
	db *gorm.DB,
	client *redis.Client,
	captcha *service.CaptchaService,
	userService *service.UserService,
	wxLoginService *service.WxLoginService,
	ipSearcher *xdb.Searcher,
	sysConfig *types.SystemConfig) *UserHandler {
	return &UserHandler{
		BaseHandler:    BaseHandler{DB: db, App: app},
		redis:          client,
		captcha:        captcha,
		userService:    userService,
		baseConfig:     &sysConfig.Base,
		wxLoginConfig:  &sysConfig.WxLogin,
		wxLoginService: wxLoginService,
		ipSearcher:     ipSearcher,
	}
}

// RegisterRoutes 注册路由
func (h *UserHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/user/")

	// 公开接口（无需登录）
	rg.POST("login", h.Login)
	rg.GET("login/qrcode", h.GetWechatQRCode)
	rg.POST("login/callback", h.WechatCallback)
	rg.GET("login/status", h.GetWechatLoginState)

	// 需要用户登录的接口
	rg.Use(middleware.UserAuthMiddleware(h.App.Config.Session.SecretKey, h.App.Redis))
	{
		rg.GET("logout", h.Logout)
		rg.GET("session", h.Session)
		rg.GET("profile", h.Profile)
		rg.POST("update/profile", h.ProfileUpdate)
		rg.POST("update/password", h.UpdatePass)
		rg.POST("update/username", h.UpdateUsername)
	}
}

// Login 用户登录（支持密码登录和验证码登录）
func (h *UserHandler) Login(c *gin.Context) {
	var data struct {
		Method   string `json:"method"` // 登录方式：password、code、wechat
		Username string `json:"username"`
		Password string `json:"password,omitempty"`
		Code     string `json:"code,omitempty"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var user model.User
	var err error
	var key string

	h.DB.Where("username", data.Username).First(&user)
	if data.Code != "" { // 验证码登录
		key = CodeStorePrefix + data.Username
		code, err := h.redis.Get(c, key).Result()
		if err != nil || code != data.Code {
			resp.ERROR(c, "验证码错误")
			return
		}

		if user.Id == 0 { // 用户不存在，创建新用户
			user, err = h.createNewUser(data.Username, "", "")
			if err != nil {
				resp.ERROR(c, err.Error())
				return
			}
		}

	} else { // 密码登录
		if user.Id == 0 {
			resp.ERROR(c, "用户不存在")
			return
		}

		// 验证密码
		genPassword := utils.GenPassword(data.Password, user.Salt)
		if genPassword != user.Password {
			resp.ERROR(c, "用户名或密码错误")
			return
		}
	}

	if !user.Enabled {
		resp.ERROR(c, "用户已禁用")
		return
	}

	// 检查用户是否已过期
	if user.ExpiredTime > 0 && user.ExpiredTime < time.Now().Unix() {
		resp.ERROR(c, "用户已过期")
		return
	}

	// 执行登录逻辑
	token, err := h.doLogin(&user, c.ClientIP())
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	// 登录成功，删除短信验证码
	if data.Code != "" {
		_ = h.redis.Del(c, key)
	}

	resp.SUCCESS(c, gin.H{
		"token":    token,
		"user_id":  user.Id,
		"username": user.Username,
	})
}

// createNewUser 创建新用户
func (h *UserHandler) createNewUser(username string, openid string, password string) (model.User, error) {
	if username == "" {
		username = fmt.Sprintf("wxuser@%d", utils.RandomNumber(8))
	}
	if password == "" {
		password = utils.RandString(8)
	}

	salt := utils.RandString(8)
	user := model.User{
		Username: username,
		Nickname: utils.GenerateNickname(),
		OpenId:   openid,
		Password: utils.GenPassword(password, salt),
		Salt:     salt,
		Enabled:  true,
		Scores:   h.baseConfig.InitScore,
		Avatar:   utils.GenerateAvatar(),
	}
	if openid != "" {
		user.Platform = "wechat"
	}

	// 检查用户名是否已存在
	var existUser model.User
	h.DB.Where("username = ?", user.Username).First(&existUser)
	if existUser.Id > 0 {
		return user, fmt.Errorf("该用户名已经被注册")
	}

	if err := h.DB.Create(&user).Error; err != nil {
		return user, fmt.Errorf("创建用户失败: %v", err)
	}

	return user, nil
}

// doLogin 执行登录操作
func (h *UserHandler) doLogin(user *model.User, ip string) (string, error) {
	// 更新最后登录时间和IP
	user.LastLoginIp = ip
	user.LastLoginAt = time.Now().Unix()
	err := h.DB.Model(user).Updates(user).Error
	if err != nil {
		return "", fmt.Errorf("failed to update user: %v", err)
	}

	// 记录登录日志
	h.DB.Create(&model.UserLoginLog{
		UserId:       user.Id,
		Username:     user.Username,
		LoginIp:      ip,
		LoginAddress: utils.Ip2Region(h.ipSearcher, ip),
	})

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"expired": time.Now().Add(time.Second * time.Duration(h.App.Config.Session.MaxAge)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(h.App.Config.Session.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	// 保存到 redis
	sessionKey := fmt.Sprintf("users/%d", user.Id)
	if _, err = h.redis.Set(context.Background(), sessionKey, tokenString, 0).Result(); err != nil {
		return "", fmt.Errorf("error with save token: %v", err)
	}

	return tokenString, nil
}

// Logout 注 销
func (h *UserHandler) Logout(c *gin.Context) {
	key := h.GetUserKey(c)
	if _, err := h.redis.Del(c, key).Result(); err != nil {
		logger.Error("error with delete session: ", err)
	}
	resp.SUCCESS(c)
}

// Session 获取/验证会话
func (h *UserHandler) Session(c *gin.Context) {
	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c, err.Error())
		return
	}

	var userVo vo.User
	err = utils.CopyObject(user, &userVo)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	// 用户 VIP 到期
	if user.ExpiredTime > 0 && user.ExpiredTime < time.Now().Unix() {
		h.DB.Model(&user).UpdateColumn("vip", false)
	}

	userVo.Id = user.Id
	resp.SUCCESS(c, userVo)

}

type userProfile struct {
	Id          uint   `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Scores      int    `json:"scores"`
	ExpiredTime int64  `json:"expired_time"`
	Vip         bool   `json:"vip"`
}

func (h *UserHandler) Profile(c *gin.Context) {
	userId := h.GetLoginUserId(c)
	if userId == 0 {
		resp.NotAuth(c)
		return
	}

	var user model.User
	h.DB.First(&user, userId)
	var profile userProfile
	err := utils.CopyObject(user, &profile)
	if err != nil {
		logger.Error("对象拷贝失败：", err.Error())
		resp.ERROR(c, "获取用户信息失败")
		return
	}

	profile.Id = user.Id
	resp.SUCCESS(c, profile)
}

func (h *UserHandler) ProfileUpdate(c *gin.Context) {
	var data userProfile
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}
	h.DB.First(&user, user.Id)
	user.Avatar = data.Avatar
	user.Nickname = data.Nickname

	res := h.DB.Select("nickname", "avatar").Updates(&user)
	if res.Error != nil {
		resp.ERROR(c, "更新用户信息失败")
		return
	}

	resp.SUCCESS(c)
}

// UpdatePass 更新密码
func (h *UserHandler) UpdatePass(c *gin.Context) {
	var data struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if len(data.Password) < 8 {
		resp.ERROR(c, "密码长度不能少于8个字符")
		return
	}

	user, err := h.GetLoginUser(c)
	if err != nil {
		resp.NotAuth(c)
		return
	}

	newPass := utils.GenPassword(data.Password, user.Salt)
	err = h.DB.Model(&user).Where("id", user.Id).UpdateColumn("password", newPass).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}

// UpdateUsername 修改用户名
func (h *UserHandler) UpdateUsername(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Code     string `json:"code"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	// 检查验证码
	key := CodeStorePrefix + data.Username
	code, err := h.redis.Get(c, key).Result()
	if err != nil || code != data.Code {
		resp.ERROR(c, "验证码错误")
		return
	}

	// 检查手机号是否被其他账号绑定
	var item model.User
	res := h.DB.Where("username", data.Username).First(&item)
	if res.Error == nil {
		resp.ERROR(c, "该账号已被使用")
		return
	}

	userId := h.GetLoginUserId(c)

	err = h.DB.Model(&item).Where("id", userId).UpdateColumn("username", data.Username).Error
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	_ = h.redis.Del(c, key) // 删除短信验证码
	resp.SUCCESS(c)
}

// GetWechatQRCode 获取微信登录二维码URL
func (h *UserHandler) GetWechatQRCode(c *gin.Context) {
	if !h.wxLoginConfig.Enabled {
		resp.ERROR(c, "微信登录功能未启用")
		return
	}

	if h.wxLoginConfig.ApiKey == "" {
		resp.ERROR(c, "微信登录服务令牌未配置")
		return
	}

	state := utils.RandString(32)
	qrCodeURL, err := h.wxLoginService.GetLoginQrCodeUrl(state)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c, gin.H{
		"url":   qrCodeURL,
		"state": state,
	})
}

// 查询微信登录状态
func (h *UserHandler) GetWechatLoginState(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		resp.ERROR(c, "参数错误")
		return
	}

	status, err := h.wxLoginService.GetLoginStatus(state)
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	if status.Status != service.LoginStatusSuccess {
		resp.SUCCESS(c, status)
		return
	}

	// 登录成功
	var user model.User
	h.DB.Where("openid = ?", status.OpenID).First(&user)
	if user.Id == 0 {
		// 创建新用户
		user, err = h.createNewUser("", status.OpenID, "")
		if err != nil {
			resp.ERROR(c, err.Error())
			return
		}
	}

	token, err := h.doLogin(&user, c.ClientIP())
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	status.Status = service.LoginStatusExpired
	h.wxLoginService.SetLoginStatus(state, *status)

	status.Status = service.LoginStatusSuccess
	status.Token = token
	resp.SUCCESS(c, status)
}

// WechatCallback 微信登录回调处理
func (h *UserHandler) WechatCallback(c *gin.Context) {
	var data struct {
		OpenID string `json:"openid"`
		State  string `json:"state"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	if data.OpenID == "" || data.State == "" {
		resp.ERROR(c, "参数错误")
		return
	}

	// 设置登录状态
	status := service.LoginStatus{
		Status: service.LoginStatusSuccess,
		OpenID: data.OpenID,
	}
	h.wxLoginService.SetLoginStatus(data.State, status)

	resp.SUCCESS(c, status)
}
