package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("oss", map[string]any{
		"key":      config.Env("ALIYUN_ACCESS_KEY_ID"),
		"secret":   config.Env("ALIYUN_ACCESS_KEY_SECRET"),
		"bucket":   config.Env("ALIYUN_BUCKET"),
		"url":      config.Env("ALIYUN_URL"),
		"endpoint": config.Env("ALIYUN_ENDPOINT"),
	})
}
