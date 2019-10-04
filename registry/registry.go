package registry

import "context"

// 服务注册插件的接口
type Registry interface {
	// 插件的名称
	Name() string

	// 初始化
	Init(ctx context.Context, opts ...Options) (err error)

	// 服务注册
	Register(ctx context.Context, service *Service) (err error)

	// 服务反注册
	Unregister(ctx context.Context, service *Service) (err error)
}
