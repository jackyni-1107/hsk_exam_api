# GoFrame 常用目标（与 gf-demo 对齐的精简版；Windows 可用 Git Bash / WSL 执行 make）

.PHONY: dao dao.all service service.all ctrl ctrl.all

# 根据数据库生成 DAO（需已配置 hack/config.yaml 且 gf 在 PATH）
dao:
	gf gen dao

dao.all:
	gf gen dao -a

service:
	gf gen service

service.all:
	gf gen service -a

ctrl:
	gf gen ctrl

ctrl.all:
	gf gen ctrl -a

.PHONY: help
help:
	@echo "Targets: dao dao.all service service.all ctrl ctrl.all cli (see hack-cli.mk)"
