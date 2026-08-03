package service

import (
	"fmt"
	"geekai/core/types"
	"geekai/store/model"
	"sync"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	db   *gorm.DB
	lock sync.Mutex
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db, lock: sync.Mutex{}}
}

// IncreaseScores 增加用户积分
func (s *UserService) IncreaseScores(userId int, score int, log model.ScoreLog) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	tx := s.db.Begin()
	err := tx.Model(&model.User{}).Where("id", userId).UpdateColumn("scores", gorm.Expr("scores + ?", score)).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var user model.User
	tx.Where("id", userId).First(&user)
	err = tx.Create(&model.ScoreLog{
		UserId:    user.Id,
		Username:  user.Username,
		Type:      log.Type,
		Amount:    score,
		Balance:   user.Scores,
		Mark:      types.ScorePlus,
		Subject:   log.Subject,
		Remark:    log.Remark,
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
func (s *UserService) DecreaseScores(userId uint, score int, log model.ScoreLog) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	tx := s.db.Begin()
	err := tx.Model(&model.User{}).Where("id", userId).UpdateColumn("scores", gorm.Expr("scores - ?", score)).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("扣减积分失败：%v", err)
	}
	var user model.User
	tx.Where("id", userId).First(&user)
	err = tx.Create(&model.ScoreLog{
		UserId:    user.Id,
		Username:  user.Username,
		Type:      log.Type,
		Amount:    score,
		Balance:   user.Scores,
		Mark:      types.ScoreSub,
		Subject:   log.Subject,
		Remark:    log.Remark,
		CreatedAt: time.Now(),
	}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("记录积分日志失败：%v", err)
	}
	tx.Commit()
	return nil
}
