package config

import clientconfig "exam/api/client/config"

type ControllerV1 struct{}

func NewV1() clientconfig.IConfig {
	return &ControllerV1{}
}
