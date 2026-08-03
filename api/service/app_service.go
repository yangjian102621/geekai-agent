package service

import (
	"errors"
	"geekai/store/model"
	"sync"

	"gorm.io/gorm"
)

type AppService struct {
	db   *gorm.DB
	lock sync.Mutex
}

func NewAppService(db *gorm.DB) *AppService {
	return &AppService{db: db, lock: sync.Mutex{}}
}

// RemoveApp 删除应用
func (s *AppService) RemoveApp(appId int, creatorId uint) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	var app model.App
	if err := s.db.Where("id = ? AND creator_id = ?", appId, creatorId).First(&app).Error; err != nil {
		return errors.New("要删除的应用不存在")
	}

	tx := s.db.Begin()

	// 删除应用
	if err := tx.Delete(&app).Error; err != nil {
		return errors.New("删除应用失败：" + err.Error())
	}

	// 删除对应应用的对话
	if err := tx.Where("app_id = ?", appId).Delete(&model.ChatItem{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除对话失败：" + err.Error())
	}

	// 删除聊天记录
	if err := tx.Where("app_id = ?", appId).Delete(&model.ChatMessage{}).Error; err != nil {
		tx.Rollback()
		return errors.New("删除聊天记录失败：" + err.Error())
	}

	tx.Commit()
	return nil
}
