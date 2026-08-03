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
	"geekai/handler"
	"geekai/service/oss"
	"geekai/store/model"
	"geekai/utils/resp"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UploadHandler struct {
	handler.BaseHandler
	uploaderManager *oss.UploaderManager
}

func NewUploadHandler(app *core.AppServer, db *gorm.DB, manager *oss.UploaderManager) *UploadHandler {
	return &UploadHandler{BaseHandler: handler.BaseHandler{DB: db, App: app}, uploaderManager: manager}
}

// RegisterRoutes 注册路由
func (h *UploadHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/upload")

	// 需要管理员登录的接口
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.POST("", h.Upload)
	}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	handler := h.uploaderManager.GetUploadHandler()
	file, err := handler.PutFile(c, "file")
	if err != nil {
		resp.ERROR(c, err.Error())
		return
	}
	userId := 0
	res := h.DB.Create(&model.File{
		UserId:    userId,
		Name:      file.Name,
		ObjKey:    file.ObjKey,
		URL:       file.URL,
		Ext:       file.Ext,
		Size:      file.Size,
		CreatedAt: time.Time{},
	})
	if res.Error != nil || res.RowsAffected == 0 {
		resp.ERROR(c, "error with update database: "+res.Error.Error())
		return
	}

	resp.SUCCESS(c, file)
}
