package main

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"context"
	"embed"
	"geekai/core"
	"geekai/core/types"
	"geekai/handler"
	"geekai/handler/admin"
	"geekai/log"
	"geekai/service"
	"geekai/service/oss"
	"geekai/service/payment"
	"geekai/service/sms"
	"geekai/store"
	"io"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var logger = log.GetLogger()

//go:embed res
var res embed.FS

// AppLifecycle 应用程序生命周期
type AppLifecycle struct {
}

// OnStart 应用程序启动时执行
func (l *AppLifecycle) OnStart(context.Context) error {
	logger.Info("AppLifecycle OnStart")
	return nil
}

// OnStop 应用程序停止时执行
func (l *AppLifecycle) OnStop(context.Context) error {
	logger.Info("AppLifecycle OnStop")
	return nil
}

func NewAppLifeCycle() *AppLifecycle {
	return &AppLifecycle{}
}

func main() {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.toml"
	}
	logger.Info("Loading config file: ", configFile)
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Panic Error:", err)
			// 打印堆栈信息
			if os.Getenv("GEEKAI_DEBUG") == "true" {
				debug.PrintStack()
			}
		}
	}()

	app := fx.New(
		// 初始化配置应用配置
		fx.Provide(func() *types.AppConfig {
			config, err := core.LoadConfig(configFile)
			if err != nil {
				logger.Fatal(err)
			}
			config.Path = configFile
			if os.Getenv("GEEKAI_DEBUG") == "true" {
				core.SaveConfig(config)
			}
			return config
		}),
		// 创建应用服务
		fx.Provide(core.NewServer),

		// 初始化数据库
		fx.Provide(store.NewGormConfig),
		fx.Provide(store.NewMysql),
		fx.Provide(store.NewRedisClient),
		// 初始化服务
		fx.Invoke(func(s *core.AppServer, client *redis.Client) {
			s.Init(client)
		}),
		fx.Provide(func(db *gorm.DB) *types.SystemConfig {
			return core.LoadSystemConfig(db)
		}),
		fx.Provide(func() embed.FS {
			return res
		}),

		// 创建 Ip2Region 查询对象
		fx.Provide(func() (*xdb.Searcher, error) {
			file, err := res.Open("res/ip2region.xdb")
			if err != nil {
				return nil, err
			}
			cBuff, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}

			return xdb.NewWithBuffer(cBuff)
		}),

		// 短信服务
		fx.Provide(sms.NewAliYunSmsService),
		fx.Provide(sms.NewBaoSmsService),
		fx.Provide(sms.NewSmsManager),

		fx.Provide(func(config *types.SystemConfig) *service.CaptchaService {
			return service.NewCaptchaService(config.Captcha)
		}),
		fx.Provide(func(config *types.SystemConfig, client *redis.Client) *service.WxLoginService {
			return service.NewWxLoginService(config.WxLogin, client)
		}),
		// 文件上传服务
		fx.Provide(oss.NewLocalStorage),
		fx.Provide(oss.NewMiniOss),
		fx.Provide(oss.NewQiNiuOss),
		fx.Provide(oss.NewAliYunOss),
		fx.Provide(oss.NewUploaderManager),
		// 邮件服务
		fx.Provide(service.NewSmtpService),
		//  用户服务
		fx.Provide(service.NewUserService),
		// 应用服务
		fx.Provide(service.NewAppService),
		fx.Provide(service.NewSnowflake),
		fx.Provide(service.NewLicenseService),
		fx.Invoke(func(s *core.AppServer, ls *service.LicenseService) {
			if os.Getenv("GEEKAI_LICENSE_SYNC") == "true" {
				ls.SyncLicense()
			}
		}),
		// coze 服务
		fx.Provide(service.NewCozeService),
		// 数据迁移
		fx.Provide(service.NewMigrationService),
		fx.Invoke(func(ms *service.MigrationService) {
			if err := ms.Migrate(); err != nil {
				logger.Errorf("数据迁移失败：%v", err)
				os.Exit(0)
			}
		}),

		fx.Provide(handler.NewUserHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.UserHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(handler.NewSmsHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.SmsHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(handler.NewCaptchaHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.CaptchaHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(admin.NewAdminHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.ManagerHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(admin.NewUserHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.UserHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(admin.NewConfigHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.ConfigHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(admin.NewUploadHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.UploadHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(handler.NewChatHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.ChatHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(admin.NewAppHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.AppHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewConfigHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.ConfigHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(handler.NewFileHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.FileHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(handler.NewAppHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.AppHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(admin.NewRedeemHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.RedeemHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(handler.NewRedeemHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.RedeemHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(admin.NewScoreLogHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.ScoreLogHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(handler.NewScoreLogHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.ScoreLogHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(handler.NewLicenseHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.LicenseHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(admin.NewProductHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.ProductHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(handler.NewProductHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.ProductHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(admin.NewOrderHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.OrderHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(handler.NewOrderHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.OrderHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(payment.NewAlipayService),
		fx.Provide(payment.NewEPayService),
		fx.Provide(payment.NewWechatService),
		fx.Provide(handler.NewPaymentHandler),
		fx.Invoke(func(s *core.AppServer, h *handler.PaymentHandler) {
			h.RegisterRoutes()
			h.StartSyncOrders()
		}),

		fx.Provide(admin.NewAppCategoryHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.AppCategoryHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(admin.NewDashboardHandler),
		fx.Invoke(func(s *core.AppServer, h *admin.DashboardHandler) {
			h.RegisterRoutes()
		}),

		// 创作者模块
		fx.Provide(service.NewCreatorService),
		fx.Provide(handler.NewCreatorHandler),
		fx.Invoke(func(h *handler.CreatorHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewCreatorAppHandler),
		fx.Invoke(func(h *handler.CreatorAppHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewAppCategoryHandler),
		fx.Invoke(func(h *handler.AppCategoryHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewWithdrawHandler),
		fx.Invoke(func(h *handler.WithdrawHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(handler.NewScoreHandler),
		fx.Invoke(func(h *handler.ScoreHandler) {
			h.RegisterRoutes()
		}),

		// 创作者后台接口
		fx.Provide(admin.NewCreatorHandler),
		fx.Invoke(func(h *admin.CreatorHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(admin.NewWithdrawHandler),
		fx.Invoke(func(h *admin.CreatorWithdrawHandler) {
			h.RegisterRoutes()
		}),
		fx.Provide(admin.NewCreatorAppHandler),
		fx.Invoke(func(h *admin.CreatorAppHandler) {
			h.RegisterRoutes()
		}),

		// 工作流模块
		fx.Provide(service.NewBailianService),
		fx.Provide(service.NewWorkflowService),
		fx.Invoke(func(s *core.AppServer, ws *service.WorkflowService) {
			go ws.StartTaskPolling()
		}),
		fx.Provide(admin.NewWorkflowHandler),
		fx.Invoke(func(h *admin.WorkflowHandler) {
			h.RegisterRoutes()
		}),

		fx.Provide(handler.NewWorkflowHandler),
		fx.Invoke(func(h *handler.WorkflowHandler) {
			h.RegisterRoutes()
		}),

		fx.Invoke(func(s *core.AppServer, db *gorm.DB, sysConfig *types.SystemConfig) {
			go func() {
				err := s.Run(db, sysConfig)
				if err != nil {
					logger.Error(err)
					os.Exit(0)
				}
			}()
		}),

		fx.Provide(NewAppLifeCycle),
		// 注册生命周期回调函数
		fx.Invoke(func(lifecycle fx.Lifecycle, lc *AppLifecycle) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return lc.OnStart(ctx)
				},
				OnStop: func(ctx context.Context) error {
					return lc.OnStop(ctx)
				},
			})
		}),
	)
	// 启动应用程序
	go func() {
		if err := app.Start(context.Background()); err != nil {
			logger.Fatal(err)
		}
	}()

	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 关闭应用程序
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		logger.Fatal(err)
	}

}
