package service

import (
	"fmt"
	"geekai/store/model"
	"geekai/utils"
	"os"
	"strings"

	"gorm.io/gorm"
)

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

// logger 已在 service/types.go 中定义，直接使用即可

// MigrationService 数据迁移服务
type MigrationService struct {
	db *gorm.DB
}

func NewMigrationService(db *gorm.DB) *MigrationService {
	return &MigrationService{db: db}
}

// Migrate 执行数据迁移
// 1. 首次运行执行所有表的 AutoMigrate
// 2. 按版本执行字段删除等操作
func (s *MigrationService) Migrate() error {
	// ==================== 第一步：首次运行执行所有表的 AutoMigrate ====================
	logger.Info("开始执行数据表迁移...")

	models := []any{
		&model.AdminUser{},
		&model.App{},
		&model.AppCategory{},
		&model.ChatItem{},
		&model.ChatMessage{},
		&model.Config{},
		&model.Creator{},
		&model.CreatorScoreLog{},
		&model.CreatorWithdraw{},
		&model.File{},
		&model.Order{},
		&model.Product{},
		&model.Redeem{},
		&model.ScoreLog{},
		&model.User{},
		&model.UserLoginLog{},
		&model.Workflow{},
		&model.WorkflowTask{},
	}

	for _, m := range models {
		if err := s.db.AutoMigrate(m); err != nil {
			logger.Errorf("数据表迁移失败 [%T]: %v", m, err)
			return fmt.Errorf("数据表迁移失败 [%T]: %w", m, err)
		}
		logger.Infof("数据表迁移成功: %T", m)
	}

	logger.Info("所有数据表迁移完成")
	if err := s.initializeAdmin(); err != nil {
		return err
	}

	// ==================== 第二步：按版本执行字段删除等操作 ====================
	if err := s.runVersionMigrations(); err != nil {
		logger.Errorf("版本迁移执行失败: %v", err)
		return err
	}

	return nil
}

func (s *MigrationService) initializeAdmin() error {
	var count int64
	if err := s.db.Model(&model.AdminUser{}).Count(&count).Error; err != nil {
		return fmt.Errorf("检查管理员账号失败: %w", err)
	}
	if count > 0 {
		return nil
	}

	username, password, err := initialAdminCredentials()
	if err != nil {
		return err
	}

	salt := utils.RandString(8)
	admin := model.AdminUser{
		Username: username,
		Password: utils.GenPassword(password, salt),
		Salt:     salt,
		Status:   true,
	}
	if err := s.db.Create(&admin).Error; err != nil {
		return fmt.Errorf("创建初始管理员失败: %w", err)
	}
	logger.Infof("已创建初始管理员账号: %s", username)
	return nil
}

func initialAdminCredentials() (string, string, error) {
	username := strings.TrimSpace(os.Getenv("GEEKAI_ADMIN_USERNAME"))
	password := os.Getenv("GEEKAI_ADMIN_PASSWORD")
	if username == "" || password == "" {
		return "", "", fmt.Errorf("系统尚无管理员，请设置 GEEKAI_ADMIN_USERNAME 和 GEEKAI_ADMIN_PASSWORD 后重新启动")
	}
	if len(username) < 3 || len(username) > 30 {
		return "", "", fmt.Errorf("GEEKAI_ADMIN_USERNAME 长度必须为 3-30 个字符")
	}
	if len(password) < 12 {
		return "", "", fmt.Errorf("GEEKAI_ADMIN_PASSWORD 长度不能少于 12 个字符")
	}
	return username, password, nil
}

// runVersionMigrations 执行版本化的迁移操作
// 每个版本的字段删除操作都在独立的代码块中
func (s *MigrationService) runVersionMigrations() error {
	// ==================== Version 1.0.0 ====================
	// 版本 1.0.5 的字段删除操作
	// 示例：删除某个表的某个字段
	// if !s.db.Migrator().HasColumn(&model.App{}, "old_field") {
	// 	s.db.Migrator().DropColumn(&model.App{}, "old_field")
	// }

	return nil
}
