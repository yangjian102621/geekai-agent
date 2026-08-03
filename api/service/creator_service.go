package service

import (
	"fmt"
	"sync"
	"time"

	"geekai/core/types"
	"geekai/store/model"

	"gorm.io/gorm"
)

type CreatorService struct {
	db   *gorm.DB
	lock sync.Mutex
}

func NewCreatorService(db *gorm.DB) *CreatorService {
	return &CreatorService{db: db, lock: sync.Mutex{}}
}

// IncreaseScores 增加用户积分
func (s *CreatorService) IncreaseScores(creatorId uint, score int, log model.CreatorScoreLog) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	tx := s.db.Begin()
	err := tx.Model(&model.Creator{}).Where("id", creatorId).UpdateColumn("scores", gorm.Expr("scores + ?", score)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var creator model.Creator
	tx.Where("id", creatorId).First(&creator)
	err = tx.Create(&model.CreatorScoreLog{
		UserId:    creator.UserId,
		AppId:     log.AppId,
		CreatorId: creator.Id,
		Type:      log.Type,
		Score:     score,
		Balance:   creator.Scores,
		Subject:   log.Subject,
		Remark:    log.Remark,
		Mark:      types.ScorePlus,
		CreatedAt: time.Now(),
	}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("记录积分日志失败：%v", err)
	}
	tx.Commit()
	return nil
}

// DecreaseScores 减少用户积分
func (s *CreatorService) DecreaseScores(creatorId uint, score int, log model.CreatorScoreLog) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	tx := s.db.Begin()
	err := tx.Model(&model.Creator{}).Where("id", creatorId).UpdateColumn("scores", gorm.Expr("scores - ?", score)).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("扣减积分失败：%v", err)
	}
	var creator model.Creator
	tx.Where("id", creatorId).First(&creator)
	err = tx.Create(&model.CreatorScoreLog{
		UserId:    creator.UserId,
		AppId:     log.AppId,
		CreatorId: creator.Id,
		Type:      log.Type,
		Score:     score,
		Balance:   creator.Scores,
		Subject:   log.Subject,
		Remark:    log.Remark,
		CreatedAt: time.Now(),
		Mark:      types.ScoreSub,
	}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("记录积分日志失败：%v", err)
	}
	tx.Commit()
	return nil
}
