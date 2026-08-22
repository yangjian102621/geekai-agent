package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"fmt"
)

type AppConfig struct {
	Path         string `toml:"-"`
	Listen       string
	Session      Session
	AdminSession Session
	ProxyURL     string
	MysqlDns     string      // mysql 连接地址
	Redis        RedisConfig // redis 连接信息
	TikaHost     string      // TiKa 服务器地址
	GeekApiHost  string      // GeekAI API 服务器地址
	AppConfigKey string      // 应用配置加密密钥
}

type RedisConfig struct {
	Host     string `json:"host"`     // 主机地址
	Port     int    `json:"port"`     // 端口
	Password string `json:"password"` // 密码
	DB       int    `json:"db"`       // 数据库
}

func (c RedisConfig) Url() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type License struct {
	Type      string        `json:"type"`
	Name      string        `json:"name"`
	License   string        `json:"license"`
	MachineId string        `json:"mid"`
	ActiveAt  int64         `json:"active_at"`
	IsActive  bool          `json:"is_active,omitempty"`
	ExpiredAt int64         `json:"expired_at"`
	Configs   LicenseConfig `json:"configs"`
}

type LicenseConfig struct {
	UserNum int  `json:"user_num"` // 用户数量
	DeCopy  bool `json:"de_copy"`  // 去版权
}

type CozeApiConfig struct {
	ApiUrl      string `json:"api_url,omitempty"`
	PrivateKey  string `json:"private_key,omitempty"`   // 授权私钥
	AppId       string `json:"app_id,omitempty"`        // 授权应用ID
	PublicKeyID string `json:"public_key_id,omitempty"` // 授权公钥ID
	SpaceId     string `json:"space_id"`                // 空间ID
}

type BailianApiConfig struct {
	ApiKey string `json:"api_key,omitempty"` // 百炼 API Key
	AppId  string `json:"app_id,omitempty"`  // 百炼应用 ID
}

type BaseConfig struct {
	Title       string `json:"title,omitempty"`        // 网站标题
	Slogan      string `json:"slogan,omitempty"`       // 网站 slogan
	AdminTitle  string `json:"admin_title,omitempty"`  // 管理后台标题
	Logo        string `json:"logo,omitempty"`         // Logo
	Copyright   string `json:"copyright,omitempty"`    // 网站版权
	InitScore   int    `json:"init_score,omitempty"`   // 新用户注册赠送积分
	DailyScore  int    `json:"daily_score,omitempty"`  // 每日签到赠送积分
	InviteScore int    `json:"invite_score,omitempty"` // 邀请新用户赠送积分

	EnabledRegister bool     `json:"enabled_register,omitempty"` // 是否开放注册
	WechatCardURL   string   `json:"wechat_card_url,omitempty"`  // 微信客服地址
	EmailWhiteList  []string `json:"email_white_list,omitempty"` // 邮箱白名单列表
	AppId           uint     `json:"app_id,omitempty"`           // 默认应用 ID
}

type SystemConfig struct {
	Base    BaseConfig    // 基础配置
	License License       // 许可证
	Coze    CozeApiConfig // Coze API 配置
	Captcha CaptchaConfig // GeekAI 行为验证码配置
	SMS     SMSConfig     // 短信发送配置
	OSS     OSSConfig     // OSS config
	Smtp    SmtpConfig    // 邮件发送配置
	Payment PaymentConfig // 支付配置
	WxLogin WxLoginConfig // 微信登录配置
}
