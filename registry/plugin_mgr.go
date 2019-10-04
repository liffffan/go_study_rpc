package registry

import (
	"context"
	"fmt"
	"sync"
)

var (
	// 定义一个全局的插件管理实例
	pluginMgr = &PluginMgr{
		plugins: make(map[string]Registry),
	}
)

// 插件管理类
type PluginMgr struct {
	// 使用 map 管理所有插件
	plugins map[string]Registry
	// 涉及到多线程操作，加了个锁
	lock sync.Mutex
}

// 注册插件到 map 里面去
func (p *PluginMgr) registerPlugin(plugin Registry) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	// 如果插件存在，直接返回错误
	_, ok := p.plugins[plugin.Name()]
	if ok {
		err = fmt.Errorf("duplicate registry plugin")
		return
	}

	p.plugins[plugin.Name()] = plugin
	return
}

// 初始化插件
func (p *PluginMgr) initPlugin(ctx context.Context, name string, opts ...Options) (registry Registry, err error) {
	// 查找对应的插件是否存在
	p.lock.Lock()
	defer p.lock.Unlock()
	plugin, ok := p.plugins[name]
	if !ok {
		err = fmt.Errorf("plugin %s not exists", name)
		return
	}

	registry = plugin
	err = plugin.Init(ctx, opts...)
	return
}

// 注册插件
func RegisterPlugin(registry Registry) (err error) {
	return pluginMgr.registerPlugin(registry)
}

// 初始化注册中心
func InitRegistry(ctx context.Context, name string, opts ...Options) (registry Registry, err error) {
	return pluginMgr.initPlugin(ctx, name)
}
