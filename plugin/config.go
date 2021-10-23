package plugin

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

// Config 插件统一配置 plugin type => { plugin name => plugin config } 。
type Config map[string]map[string]yaml.Node

// Setup 通过配置生成并装载具体插件。
func (c Config) Setup() error {
	var (
		pluginChan  = make(chan Info, MaxPluginSize) // 初始化插件队列，方便后面按顺序逐个加载插件
		setupStatus = make(map[string]bool)          // 插件初始化状态，plugin key => true初始化完成 false未初始化
	)

	// 从框架配置文件中逐个取出插件并放进channel队列中
	for typ, factories := range c {
		for name, cfg := range factories {
			factory := Get(typ, name)
			if factory == nil {
				return fmt.Errorf("plugin %s:%s no registered or imported, do not configure", typ, name)
			}
			p := Info{
				factory: factory,
				typ:     typ,
				name:    name,
				cfg:     cfg,
			}
			select {
			case pluginChan <- p:
			default:
				return fmt.Errorf("plugin number exceed max limit:%d", len(pluginChan))
			}
			setupStatus[p.Key()] = false
		}
	}

	// 从channel队列中取出插件并初始化
	num := len(pluginChan)
	for num > 0 {
		for i := 0; i < num; i++ {
			p := <-pluginChan
			if err := p.Setup(); err != nil {
				return err
			}
			setupStatus[p.Key()] = true
		}
		if len(pluginChan) == num { // 循环依赖导致无插件可以初始化，返回失败
			return fmt.Errorf("cycle depends, not plugin is setup")
		}
		num = len(pluginChan)
	}

	// 发出插件初始化完成通知，个别业务逻辑需要依赖插件完成才能继续往下执行
	select {
	case <-done: // 已经close过了
	default:
		close(done)
	}
	return nil
}
