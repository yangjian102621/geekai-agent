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
	"geekai/utils"
	"geekai/utils/resp"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	handler.BaseHandler
	redis *redis.Client
}

func NewDashboardHandler(app *core.AppServer, db *gorm.DB, client *redis.Client) *DashboardHandler {
	return &DashboardHandler{
		BaseHandler: handler.BaseHandler{App: app, DB: db},
		redis:       client,
	}
}

// RegisterRoutes 注册路由
func (h *DashboardHandler) RegisterRoutes() {
	rg := h.App.Engine.Group("/api/admin/dashboard")

	// 需要管理员登录的接口
	rg.Use(middleware.AdminAuthMiddleware(h.App.Config.AdminSession.SecretKey, h.App.Redis))
	{
		rg.GET("/stats", h.GetStats)
		rg.GET("/trends", h.GetTrends)
		rg.GET("/recent", h.GetRecent)
		rg.GET("/system", h.GetSystem)
	}
}

// DashboardStats Dashboard统计数据
type DashboardStats struct {
	UserTotal      int64   `json:"user_total"`      // 用户总数
	UserTodayNew   int64   `json:"user_today_new"`  // 今日新增用户
	UserActive     int64   `json:"user_active"`     // 活跃用户数(7天内)
	AppTotal       int64   `json:"app_total"`       // 应用总数
	AppEnabled     int64   `json:"app_enabled"`     // 启用应用数
	ChatTotal      int64   `json:"chat_total"`      // 对话总数
	ChatToday      int64   `json:"chat_today"`      // 今日对话数
	OrderTotal     int64   `json:"order_total"`     // 订单总数
	OrderToday     int64   `json:"order_today"`     // 今日订单数
	RevenueTotal   float64 `json:"revenue_total"`   // 总收入
	RevenueToday   float64 `json:"revenue_today"`   // 今日收入
	ScoreConsumed  int64   `json:"score_consumed"`  // 积分消费总数
	ScoreRecharged int64   `json:"score_recharged"` // 积分充值总数
	CreatorTotal   int64   `json:"creator_total"`   // 创作者总数
	CreatorActive  int64   `json:"creator_active"`  // 活跃创作者数
}

// TrendData 趋势数据
type TrendData struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// DashboardTrends Dashboard趋势数据
type DashboardTrends struct {
	UserTrend    []TrendData `json:"user_trend"`    // 用户增长趋势
	RevenueTrend []TrendData `json:"revenue_trend"` // 收入趋势
	ChatTrend    []TrendData `json:"chat_trend"`    // 对话量趋势
}

// RecentOrder 最近订单
type RecentOrder struct {
	Id       uint    `json:"id"`
	OrderNo  string  `json:"order_no"`
	Username string  `json:"username"`
	Subject  string  `json:"subject"`
	Amount   float64 `json:"amount"`
	Status   string  `json:"status"`
	PayTime  int64   `json:"pay_time"`
}

// RecentUser 最近用户
type RecentUser struct {
	Id        uint   `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Platform  string `json:"platform"`
	CreatedAt int64  `json:"created_at"`
}

// HotApp 热门应用
type HotApp struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	UseCount int    `json:"use_count"`
	Score    int    `json:"score"`
}

// DashboardRecent Dashboard最近数据
type DashboardRecent struct {
	RecentOrders []RecentOrder `json:"recent_orders"` // 最近订单
	RecentUsers  []RecentUser  `json:"recent_users"`  // 最近用户
	HotApps      []HotApp      `json:"hot_apps"`      // 热门应用
}

// SystemStatus 系统状态
type SystemStatus struct {
	DatabaseStatus bool   `json:"database_status"` // 数据库状态
	RedisStatus    bool   `json:"redis_status"`    // Redis状态
	SystemLoad     string `json:"system_load"`     // 系统负载
	MemoryUsage    string `json:"memory_usage"`    // 内存使用
	DiskUsage      string `json:"disk_usage"`      // 磁盘使用
}

const (
	CacheKeyStats  = "geekai:agent:dashboard:stats"
	CacheKeyTrends = "geekai:agent:dashboard:trends"
	CacheKeyRecent = "geekai:agent:dashboard:recent"
	CacheKeySystem = "geekai:agent:dashboard:system"
	CacheKeyExpire = 3 * time.Minute
)

// GetStats 获取统计数据
func (h *DashboardHandler) GetStats(c *gin.Context) {
	cacheKey := CacheKeyStats

	// 尝试从Redis获取缓存数据
	cached, err := h.redis.Get(c, cacheKey).Result()
	if err == nil {
		var stats DashboardStats
		if utils.JsonDecode(cached, &stats) == nil {
			resp.SUCCESS(c, stats)
			return
		}
	}

	var stats DashboardStats
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	sevenDaysAgo := today.AddDate(0, 0, -7)

	// 用户统计
	h.DB.Model(&model.User{}).Count(&stats.UserTotal)
	h.DB.Model(&model.User{}).Where("created_at >= ?", today).Count(&stats.UserTodayNew)
	h.DB.Model(&model.User{}).Where("last_login_at >= ?", sevenDaysAgo.Unix()).Count(&stats.UserActive)

	// 应用统计
	h.DB.Model(&model.App{}).Count(&stats.AppTotal)
	h.DB.Model(&model.App{}).Where("enabled = ?", true).Count(&stats.AppEnabled)

	// 对话统计
	h.DB.Model(&model.ChatMessage{}).Count(&stats.ChatTotal)
	h.DB.Model(&model.ChatMessage{}).Where("created_at >= ?", today).Count(&stats.ChatToday)

	// 订单统计
	h.DB.Model(&model.Order{}).Count(&stats.OrderTotal)
	h.DB.Model(&model.Order{}).Where("created_at >= ?", today).Count(&stats.OrderToday)

	// 收入统计
	h.DB.Model(&model.Order{}).Where("status = ?", types.OrderPaidSuccess).Select("COALESCE(SUM(amount), 0)").Scan(&stats.RevenueTotal)
	h.DB.Model(&model.Order{}).Where("status = ? AND created_at >= ?", types.OrderPaidSuccess, today).Select("COALESCE(SUM(amount), 0)").Scan(&stats.RevenueToday)

	// 积分统计
	h.DB.Model(&model.ScoreLog{}).Where("mark = ?", types.ScoreSub).Select("COALESCE(SUM(amount), 0)").Scan(&stats.ScoreConsumed)
	h.DB.Model(&model.ScoreLog{}).Where("mark = ?", types.ScorePlus).Select("COALESCE(SUM(amount), 0)").Scan(&stats.ScoreRecharged)

	// 创作者统计（如果有创作者模块）
	if h.DB.Migrator().HasTable("geekai_creators") {
		h.DB.Table("geekai_creators").Count(&stats.CreatorTotal)
		h.DB.Table("geekai_creators").Where("enabled = ? AND `check` = ?", true, 1).Count(&stats.CreatorActive)
	}

	// 缓存结果3分钟
	if statsJson := utils.JsonEncode(stats); statsJson != "" {
		h.redis.Set(c, cacheKey, statsJson, CacheKeyExpire)
	}

	resp.SUCCESS(c, stats)
}

// GetTrends 获取趋势数据
func (h *DashboardHandler) GetTrends(c *gin.Context) {
	days := h.GetInt(c, "days", 7) // 默认7天
	// 最多30天
	days = int(math.Min(float64(days), 30))

	cacheKey := fmt.Sprintf("%s:%d", CacheKeyTrends, days)

	// 尝试从Redis获取缓存数据
	cached, err := h.redis.Get(c, cacheKey).Result()
	if err == nil {
		var trends DashboardTrends
		if utils.JsonDecode(cached, &trends) == nil {
			resp.SUCCESS(c, trends)
			return
		}
	}

	var trends DashboardTrends
	now := time.Now()

	// 生成日期数组
	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("01-02")
		dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		dayEnd := dayStart.Add(24 * time.Hour)

		// 用户增长趋势
		var userCount int64
		h.DB.Model(&model.User{}).Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).Count(&userCount)
		trends.UserTrend = append(trends.UserTrend, TrendData{
			Date:  dateStr,
			Value: float64(userCount),
		})

		// 收入趋势
		var revenue float64
		h.DB.Model(&model.Order{}).Where("status = ? AND created_at >= ? AND created_at < ?", types.OrderPaidSuccess, dayStart, dayEnd).Select("COALESCE(SUM(amount), 0)").Scan(&revenue)
		trends.RevenueTrend = append(trends.RevenueTrend, TrendData{
			Date:  dateStr,
			Value: revenue,
		})

		// 对话量趋势
		var chatCount int64
		h.DB.Model(&model.ChatMessage{}).Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).Count(&chatCount)
		trends.ChatTrend = append(trends.ChatTrend, TrendData{
			Date:  dateStr,
			Value: float64(chatCount),
		})
	}

	// 缓存结果3分钟
	if trendsJson := utils.JsonEncode(trends); trendsJson != "" {
		h.redis.Set(c, cacheKey, trendsJson, CacheKeyExpire)
	}

	resp.SUCCESS(c, trends)
}

// GetRecent 获取最近数据
func (h *DashboardHandler) GetRecent(c *gin.Context) {
	cacheKey := CacheKeyRecent

	// 尝试从Redis获取缓存数据
	cached, err := h.redis.Get(c, cacheKey).Result()
	if err == nil {
		var recent DashboardRecent
		if utils.JsonDecode(cached, &recent) == nil {
			resp.SUCCESS(c, recent)
			return
		}
	}

	var recent DashboardRecent

	// 最近10个订单
	var orders []model.Order
	h.DB.Where("status = ?", types.OrderPaidSuccess).Order("created_at DESC").Limit(10).Find(&orders)
	for _, order := range orders {
		statusName := "已完成"
		if order.Status == types.OrderNotPaid {
			statusName = "待支付"
		} else if order.Status == types.OrderPaidFailed {
			statusName = "支付失败"
		}

		recent.RecentOrders = append(recent.RecentOrders, RecentOrder{
			Id:       order.Id,
			OrderNo:  order.OrderNo,
			Username: order.Username,
			Subject:  order.Subject,
			Amount:   order.Amount,
			Status:   statusName,
			PayTime:  order.PayTime,
		})
	}

	// 最近10个用户
	var users []model.User
	h.DB.Order("created_at DESC").Limit(10).Find(&users)
	for _, user := range users {
		recent.RecentUsers = append(recent.RecentUsers, RecentUser{
			Id:        user.Id,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			Platform:  user.Platform,
			CreatedAt: user.CreatedAt.Unix(),
		})
	}

	// 热门应用Top10
	var apps []model.App
	h.DB.Where("enabled = ?", true).Order("use_count DESC").Limit(10).Find(&apps)
	for _, app := range apps {
		recent.HotApps = append(recent.HotApps, HotApp{
			Id:       app.Id,
			Name:     app.Name,
			Icon:     app.Icon,
			UseCount: app.UseCount,
			Score:    app.Score,
		})
	}

	// 缓存结果1分钟
	if recentJson := utils.JsonEncode(recent); recentJson != "" {
		h.redis.Set(c, cacheKey, recentJson, CacheKeyExpire)
	}

	resp.SUCCESS(c, recent)
}

// GetSystem 获取系统状态
func (h *DashboardHandler) GetSystem(c *gin.Context) {
	cacheKey := CacheKeySystem

	// 尝试从Redis获取缓存数据
	cached, err := h.redis.Get(c, cacheKey).Result()
	if err == nil {
		var status SystemStatus
		if utils.JsonDecode(cached, &status) == nil {
			resp.SUCCESS(c, status)
			return
		}
	}

	var status SystemStatus

	// 检查数据库状态
	sqlDB, err := h.DB.DB()
	if err != nil {
		status.DatabaseStatus = false
	} else {
		err = sqlDB.Ping()
		status.DatabaseStatus = err == nil
	}

	// 检查Redis状态
	_, err = h.redis.Ping(c).Result()
	status.RedisStatus = err == nil

	// 系统信息（简化版本，实际应用中可以使用系统监控库获取详细信息）
	load, _ := load.Avg()
	mem, _ := mem.VirtualMemory()
	disk, _ := disk.Usage("/")
	status.SystemLoad = fmt.Sprintf("%.2f, %.2f, %.2f", load.Load1, load.Load5, load.Load15)
	status.MemoryUsage = fmt.Sprintf("%.2f", mem.UsedPercent)
	status.DiskUsage = fmt.Sprintf("%.2f", disk.UsedPercent)

	// 缓存结果1分钟
	if statusJson := utils.JsonEncode(status); statusJson != "" {
		h.redis.Set(c, cacheKey, statusJson, CacheKeyExpire)
	}

	resp.SUCCESS(c, status)
}
