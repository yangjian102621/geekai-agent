package core

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
	"geekai/core/middleware"
	"geekai/core/types"
	"geekai/utils"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/imroc/req/v3"
	"github.com/shirou/gopsutil/host"
	"gorm.io/gorm"
)

type AppServer struct {
	Debug  bool
	Config *types.AppConfig
	Engine *gin.Engine
	Redis  *redis.Client
}

// AuthConfig 定义授权配置
type AuthConfig struct {
	ExactPaths  map[string]bool // 精确匹配的路径
	PrefixPaths map[string]bool // 前缀匹配的路径
	ParamPaths  map[string]bool // 参数化路径匹配，如 /api/creator/:username
}

var authConfig = &AuthConfig{
	ExactPaths: map[string]bool{
		"/api/user/login":        false,
		"/api/user/logout":       false,
		"/api/app/list":          false,
		"/api/admin/login":       false,
		"/api/admin/logout":      false,
		"/api/user/register":     false,
		"/api/license/get":       false,
		"/api/app/hot":           false,
		"/api/app/category/list": false,
	},
	PrefixPaths: map[string]bool{
		"/api/test/":       false,
		"/api/config/":     false,
		"/api/sms/":        false,
		"/api/captcha/":    false,
		"/api/user/login/": false,
	},
	ParamPaths: map[string]bool{},
}

func NewServer(appConfig *types.AppConfig, client *redis.Client) *AppServer {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	return &AppServer{
		Debug:  false,
		Config: appConfig,
		Engine: gin.Default(),
		Redis:  client,
	}
}

func (s *AppServer) Init(client *redis.Client) {
	s.Engine.Use(middleware.ParameterHandlerMiddleware())
	s.Engine.Use(errorHandler)

	// 添加静态资源访问
	s.Engine.Static("/static", "./static")
}

func (s *AppServer) Run(db *gorm.DB, sysConfig *types.SystemConfig) error {
	// 统计安装信息
	go func() {
		info, err := host.Info()
		if err == nil {
			apiURL := fmt.Sprintf("%s/%s", sysConfig.Captcha.ApiURL, "api/installs/push")
			timestamp := time.Now().Unix()
			product := "geek-agent"
			signStr := fmt.Sprintf("%s#%s#%d", product, info.HostID, timestamp)
			sign := utils.Sha256(signStr)
			r, err := req.C().R().SetBody(map[string]any{
				"product":   product,
				"device_id": info.HostID,
				"timestamp": timestamp,
				"sign":      sign,
			}).Post(apiURL)
			if err != nil {
				logger.Errorf("register install info failed: %v", err)
			} else {
				logger.Debugf("register install info success: %v", r.String())
			}
		}
	}()

	logger.Infof("http://%s", s.Config.Listen)
	return s.Engine.Run(s.Config.Listen)
}

// 全局异常处理
func errorHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Handler Panic: %v", r)
			debug.PrintStack()
			c.JSON(http.StatusBadRequest, types.BizVo{Code: types.Failed, Message: types.ErrorMsg})
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}
