package mock

// Mock 预留：客户端模拟卷等业务服务入口（当前由 controller 直连 DAO）。
type sMock struct{}

var Mock = new(sMock)
