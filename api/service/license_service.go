package service

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import (
	"encoding/json"
	"fmt"
	"geekai/core/types"
	"geekai/utils"
	"time"

	"geekai/store/model"

	"github.com/imroc/req/v3"
	"github.com/shirou/gopsutil/host"
	"gorm.io/gorm"
)

type LicenseService struct {
	apiURL    string
	license   *types.License
	machineId string
	db        *gorm.DB
}

func NewLicenseService(sysConfig *types.SystemConfig, db *gorm.DB, appConfig *types.AppConfig) *LicenseService {
	var machineId string
	info, err := host.Info()
	if err == nil {
		machineId = info.HostID
	}

	return &LicenseService{
		apiURL:    appConfig.GeekApiHost,
		license:   &sysConfig.License,
		machineId: machineId,
		db:        db,
	}
}

// ActiveLicense 激活 License
func (s *LicenseService) ActiveLicense(license string) (*types.License, error) {
	var res struct {
		Code    types.BizCode `json:"code"`
		Message string        `json:"message"`
		Data    types.License `json:"data"`
	}
	apiURL := fmt.Sprintf("%s/%s", s.apiURL, "api/license/active")
	response, err := req.C().R().
		SetBody(map[string]string{"license": license, "machine_id": s.machineId}).
		SetSuccessResult(&res).Post(apiURL)
	if err != nil {
		return nil, fmt.Errorf("发送激活请求失败: %v", err)
	}

	if response.IsErrorState() {
		return nil, fmt.Errorf("发送激活请求失败：%v", response.Status)
	}

	if res.Code != types.Success {
		return nil, fmt.Errorf("激活失败：%v", res.Message)
	}

	s.license = &types.License{
		License:   license,
		MachineId: s.machineId,
		IsActive:  true,
		ActiveAt:  res.Data.ActiveAt,
		ExpiredAt: res.Data.ExpiredAt,
		Configs:   res.Data.Configs,
	}

	return s.license, nil
}

// SyncLicense 定期同步 License
func (s *LicenseService) SyncLicense() {
	go func() {
		// 申请免费许可证
		if s.license.License == "" {
			_, err := s.ApplyFreeLicense()
			if err != nil {
				logger.Debugf("failed to apply free license: %v", err)
			}
		}

		// 同步许可证书
		for {
			if s.license.License == "" {
				time.Sleep(time.Second * 10)
				continue
			}

			err := s.fetchLicense()
			if err != nil {
				logger.Debugf("同步许可证书失败: %v", err)
				s.license.IsActive = false
			} else {
				s.license.IsActive = true
			}

			time.Sleep(time.Second * 60)
		}
	}()
}

// ApplyFreeLicense 申请免费License
func (s *LicenseService) ApplyFreeLicense() (*types.License, error) {
	info, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get host info: %v", err)
	}

	apiURL := fmt.Sprintf("%s/%s", s.apiURL, "api/license/apply")
	timestamp := time.Now().Unix()
	body := map[string]any{
		"type":       "free",
		"product":    "geek-agent",
		"machine_id": info.HostID,
		"timestamp":  timestamp,
	}
	signStr := fmt.Sprintf("%s@#%s@#%s@#%d@#", body["type"], body["product"], body["machine_id"], body["timestamp"])
	sign := utils.Sha256(signStr)
	body["sign"] = sign
	resp, err := req.C().R().SetBody(body).Post(apiURL)
	if err != nil {
		return nil, fmt.Errorf("apply license failed: %v", err)
	}

	var res struct {
		Code    types.BizCode `json:"code"`
		Message string        `json:"message"`
		Data    types.License `json:"data"`
	}
	err = json.Unmarshal([]byte(resp.String()), &res)
	if err != nil {
		return nil, fmt.Errorf("unmarshal license response failed: %v", err)
	}

	if res.Code != types.Success {
		return nil, fmt.Errorf("apply license failed: %v", res.Message)
	}

	// 检查是否过期
	if time.Now().Unix() > res.Data.ExpiredAt {
		res.Data.IsActive = false
	}

	s.license = &res.Data
	return s.license, nil
}

func (s *LicenseService) fetchLicense() error {
	var res struct {
		Code    types.BizCode `json:"code"`
		Message string        `json:"message"`
		Data    types.License `json:"data"`
	}

	apiURL := fmt.Sprintf("%s/%s", s.apiURL, "api/license/check")
	logger.Debugf("apiURL: %s", apiURL)
	response, err := req.C().R().
		SetBody(map[string]string{"license": s.license.License, "machine_id": s.machineId}).
		SetSuccessResult(&res).Post(apiURL)
	if err != nil {
		return fmt.Errorf("发送同步许可证书请求失败: %v", err)
	}
	if response.IsErrorState() {
		return fmt.Errorf("同步许可证书失败：%v", response.Status)
	}
	if res.Code != types.Success {
		return fmt.Errorf("同步许可证书失败：%v", res.Message)
	}

	if res.Data.ExpiredAt < time.Now().Unix() {
		logger.Warn("许可证书已经过期")
		s.license.IsActive = false
	} else {
		s.license.Configs = res.Data.Configs
		s.license.ExpiredAt = res.Data.ExpiredAt
		s.license.Type = res.Data.Type
		s.license.Name = res.Data.Name
		// 更新数据库
		s.db.Model(&model.Config{}).Where("name = ?", "license").Update("value", utils.JsonEncode(s.license))
	}
	logger.Debugf("同步许可证书成功: %+v", res.Data)
	return nil
}

// GetLicense 获取许可信息
func (s *LicenseService) GetLicense() *types.License {
	return s.license
}
