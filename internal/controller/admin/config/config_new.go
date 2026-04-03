package config

import "exam/api/admin/config"

type ControllerV1 struct{}

func NewV1() config.IConfig {
	return &ControllerV1{}
}
