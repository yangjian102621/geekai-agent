package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"geekai/core/types"
	"geekai/store/vo"
	"io"
	"strings"
)

const encryptedAppConfigPrefix = "enc:v1:"

// AppConfigService 负责应用配置中的密钥加密、解密和脱敏。
type AppConfigService struct {
	key []byte
}

func NewAppConfigService(config *types.AppConfig) (*AppConfigService, error) {
	keyMaterial := strings.TrimSpace(config.AppConfigKey)
	if len(keyMaterial) < 16 {
		return nil, errors.New("应用配置加密密钥未配置或长度不足，请设置 GEEKAI_APP_CONFIG_KEY（至少 16 个字符）")
	}

	key := sha256.Sum256([]byte(keyMaterial))
	return &AppConfigService{key: key[:]}, nil
}

// Encode 将应用配置序列化并加密后保存到 geekai_apps.configs。
func (s *AppConfigService) Encode(config vo.AppConfig) (string, error) {
	raw, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("序列化应用配置失败: %w", err)
	}
	return s.EncryptRaw(string(raw))
}

// Decode 读取应用配置。为兼容升级，未加密的旧数据仍可读取，后续会在启动迁移时加密。
func (s *AppConfigService) Decode(raw string, dest any) error {
	plain, err := s.DecryptRaw(raw)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(plain), dest); err != nil {
		return fmt.Errorf("解析应用配置失败: %w", err)
	}
	return nil
}

func (s *AppConfigService) EncryptRaw(raw string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", fmt.Errorf("初始化应用配置加密失败: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("初始化应用配置加密失败: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("生成应用配置加密随机数失败: %w", err)
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(raw), nil)
	payload := append(nonce, ciphertext...)
	return encryptedAppConfigPrefix + base64.RawStdEncoding.EncodeToString(payload), nil
}

func (s *AppConfigService) DecryptRaw(raw string) (string, error) {
	if !strings.HasPrefix(raw, encryptedAppConfigPrefix) {
		return raw, nil
	}

	payload, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(raw, encryptedAppConfigPrefix))
	if err != nil {
		return "", fmt.Errorf("解析应用配置密文失败: %w", err)
	}
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", fmt.Errorf("初始化应用配置解密失败: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("初始化应用配置解密失败: %w", err)
	}
	if len(payload) < gcm.NonceSize() {
		return "", errors.New("应用配置密文格式无效")
	}

	nonce, ciphertext := payload[:gcm.NonceSize()], payload[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("应用配置解密失败，请检查 GEEKAI_APP_CONFIG_KEY 是否正确")
	}
	return string(plain), nil
}

// Mask 清除所有会被当作凭据使用的字段，供 HTTP 响应使用。
func (s *AppConfigService) Mask(config vo.AppConfig) vo.AppConfig {
	config.Token = ""
	config.PrivateKey = ""
	config.BailianApiKey = ""
	return config
}

// MergeSecrets 在编辑应用时保留未重新输入的旧密钥，避免前端脱敏后覆盖原密钥。
func (s *AppConfigService) MergeSecrets(current, incoming vo.AppConfig) vo.AppConfig {
	if incoming.Token == "" {
		incoming.Token = current.Token
	}
	if incoming.PrivateKey == "" {
		incoming.PrivateKey = current.PrivateKey
	}
	if incoming.BailianApiKey == "" {
		incoming.BailianApiKey = current.BailianApiKey
	}
	return incoming
}
