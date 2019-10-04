package registry

import "time"

// 选项模式
type Options struct {
	// 有两个选项，一个是服务注册中心的地址
	Addrs []string
	// 一个是与注册中心交互的超时时间
	Timeout time.Duration
}

type Option func(opts *Options)

// 初始化 Timeout
func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.Timeout = timeout
	}
}

// 初始化注册中心地址
func WithAddrs(addrs []string) Option {
	return func(opts *Options) {
		opts.Addrs = addrs
	}
}
