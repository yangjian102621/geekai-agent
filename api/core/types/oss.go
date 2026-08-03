package types

// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type OSSConfig struct {
	Active string             `json:"active"`
	Minio  MiniOssConfig      `json:"minio"`
	QiNiu  QiNiuOssConfig     `json:"qiniu"`
	AliYun AliYunOssConfig    `json:"aliyun"`
	Local  LocalStorageConfig `json:"local"`
}
type MiniOssConfig struct {
	Endpoint     string `json:"endpoint"`
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
	Bucket       string `json:"bucket"`
	UseSSL       bool   `json:"use_ssl"`
	Domain       string `json:"domain"`
}

func (c *MiniOssConfig) Equal(other *MiniOssConfig) bool {
	return c.Endpoint == other.Endpoint &&
		c.AccessKey == other.AccessKey &&
		c.AccessSecret == other.AccessSecret &&
		c.Bucket == other.Bucket &&
		c.UseSSL == other.UseSSL &&
		c.Domain == other.Domain
}

type QiNiuOssConfig struct {
	Zone         string `json:"zone"`
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
	Bucket       string `json:"bucket"`
	Domain       string `json:"domain"`
}

func (c *QiNiuOssConfig) Equal(other *QiNiuOssConfig) bool {
	return c.Zone == other.Zone &&
		c.AccessKey == other.AccessKey &&
		c.AccessSecret == other.AccessSecret &&
		c.Bucket == other.Bucket &&
		c.Domain == other.Domain
}

type AliYunOssConfig struct {
	Endpoint     string `json:"endpoint"`
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
	Bucket       string `json:"bucket"`
	Domain       string `json:"domain"`
}

func (c *AliYunOssConfig) Equal(other *AliYunOssConfig) bool {
	return c.Endpoint == other.Endpoint &&
		c.AccessKey == other.AccessKey &&
		c.AccessSecret == other.AccessSecret &&
		c.Bucket == other.Bucket &&
		c.Domain == other.Domain
}

type LocalStorageConfig struct {
	BasePath string `json:"base_path"`
	BaseURL  string `json:"base_url"`
}

func (c *LocalStorageConfig) Equal(other *LocalStorageConfig) bool {
	return c.BasePath == other.BasePath &&
		c.BaseURL == other.BaseURL
}
