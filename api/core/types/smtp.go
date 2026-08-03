package types

type SmtpConfig struct {
	UseTls   bool   `json:"use_tls"`  // 是否使用 TLS 发送
	Host     string `json:"host"`     // 邮件服务器地址
	Port     int    `json:"port"`     // 邮件服务器端口
	AppName  string `json:"app_name"` // 应用名称
	From     string `json:"from"`     // 发件人邮箱地址
	Password string `json:"password"` // 发件人邮箱密码
}

func (s *SmtpConfig) Equal(other *SmtpConfig) bool {
	return s.UseTls == other.UseTls &&
		s.Host == other.Host &&
		s.Port == other.Port &&
		s.AppName == other.AppName &&
		s.From == other.From &&
		s.Password == other.Password
}
